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

type WrapAsImage struct {
	SourceDir         string
	TargetDir         string
	ParentImage       string
	NewRepositoryName string
}

func (img WrapAsImage) DryRun() {
	log.WithFields(log.Fields{"dry": true}).Info("WRAP")
}

func (img WrapAsImage) WithDockerClient(c *docker.Client) error {
	imageId := utils.RandomIdentifierOfLength(64)

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

	err = pack_it_up(img.SourceDir, layerFile, img.TargetDir)
	if err != nil {
		return err
	}
	layerFile.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	log.WithFields(log.Fields{"image_id": imageId, "repository": img.NewRepositoryName}).Info("Wrapping up")

	go func() {
		defer wg.Done()
		defer uploadWriter.Close()
		writeUploadBallInto(uploadWriter, layerBallName, img.NewRepositoryName, parentImageConfig.ID, imageId)
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

func (img WrapAsImage) AsShellCommand() string {
	return fmt.Sprintf("echo NOT IMPLEMENTED YET\n")
}

func writeUploadBallInto(w io.Writer, layerBallName string, newRepositoryName string, parentImageId string, imageId string) error {
	uploadBall := tar.NewWriter(w)
	defer uploadBall.Close()

	repositoriesFileHeader, repoFile := repositoriesFile(newRepositoryName, "latest", imageId)
	uploadBall.WriteHeader(&repositoriesFileHeader)
	uploadBall.Write(repoFile)

	imageDirHeader := imageDir(imageId)
	uploadBall.WriteHeader(&imageDirHeader)

	versionFileHeader, versionFileContents := versionFile(imageId)
	uploadBall.WriteHeader(&versionFileHeader)
	uploadBall.Write(versionFileContents)

	configFileHeader, configFileContents := imageConfigFile(parentImageId, imageId)
	uploadBall.WriteHeader(&configFileHeader)
	uploadBall.Write(configFileContents)

	info, err := os.Stat(layerBallName)
	if err != nil {
		return err
	}
	layerBallHeader := tar.Header{
		Name:     path.Join(imageId, "layer.tar"),
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
