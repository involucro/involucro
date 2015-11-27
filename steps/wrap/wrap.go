package wrap

import (
	"archive/tar"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	utils "github.com/thriqon/involucro/lib"
	"io"
	"os"
	"path"
	"sync"
)

// AsImage represents packing up a directory into
// an image derived from another.
type AsImage struct {
	SourceDir         string
	TargetDir         string
	ParentImage       string
	NewRepositoryName string
}

// DryRun runs this task without doing anything, but logging an indication of
// what would have been done
func (img AsImage) DryRun() {
	log.WithFields(log.Fields{"dry": true}).Info("WRAP")
}

// WithDockerClient executes the task on the given Docker instance
func (img AsImage) WithDockerClient(c *docker.Client) error {
	imageID := utils.RandomIdentifierOfLength(64)

	parentImageConfig, err := c.InspectImage(img.ParentImage)
	if err != nil {
		return err
	}

	uploadReader, uploadWriter := io.Pipe()

	layerBallName := randomTarballFileName()

	layerFile, err := os.Create(layerBallName)
	if err != nil {
		return err
	}
	defer os.Remove(layerBallName)
	defer layerFile.Close()

	log.WithFields(log.Fields{"layerBallName": layerBallName}).Debug("Packing")

	err = packItUp(img.SourceDir, layerFile, img.TargetDir)
	if err != nil {
		return err
	}
	layerFile.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	log.WithFields(log.Fields{"image_id": imageID, "repository": img.NewRepositoryName}).Info("Wrapping up")

	go func() {
		defer wg.Done()
		defer uploadWriter.Close()
		writeUploadBallInto(uploadWriter, layerBallName, img.NewRepositoryName, parentImageConfig.ID, imageID)
		return
	}()

	go func(r io.Reader) {
		defer wg.Done()
		err := c.LoadImage(docker.LoadImageOptions{InputStream: r})
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("Failed loading image into Docker")
		} else {
			log.Debug("Upload finished")
		}
	}(uploadReader)

	wg.Wait()
	log.Info("DONE")

	return nil
}

func writeUploadBallInto(w io.Writer, layerBallName string, newRepositoryName string, parentImageID string, imageID string) error {
	uploadBall := tar.NewWriter(w)
	defer uploadBall.Close()

	repositoriesFileHeader, repoFile := repositoriesFile(newRepositoryName, "latest", imageID)
	uploadBall.WriteHeader(&repositoriesFileHeader)
	uploadBall.Write(repoFile)

	imageDirHeader := imageDir(imageID)
	uploadBall.WriteHeader(&imageDirHeader)

	versionFileHeader, versionFileContents := versionFile(imageID)
	uploadBall.WriteHeader(&versionFileHeader)
	uploadBall.Write(versionFileContents)

	configFileHeader, configFileContents := imageConfigFile(parentImageID, imageID)
	uploadBall.WriteHeader(&configFileHeader)
	uploadBall.Write(configFileContents)

	info, err := os.Stat(layerBallName)
	if err != nil {
		return err
	}
	layerBallHeader := tar.Header{
		Name:     path.Join(imageID, "layer.tar"),
		Typeflag: tar.TypeReg,
		Size:     info.Size(),
	}
	uploadBall.WriteHeader(&layerBallHeader)
	layerBallFile, err := os.Open(layerBallName)
	if err != nil {
		return err
	}
	defer layerBallFile.Close()

	_, err = io.Copy(uploadBall, layerBallFile)

	log.Debug("Pipe finished")

	return err
}

// AsShellCommandOn prints sh compatible commands into the given writer, that
// accomplish the funciontality encoded in this step
func (img AsImage) AsShellCommandOn(w io.Writer) {
	fmt.Fprintf(w, "echo NOT IMPLEMENTED YET\n")
}
