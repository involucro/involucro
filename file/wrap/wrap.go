package wrap

import (
	"archive/tar"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/file/pull"
	utils "github.com/thriqon/involucro/lib"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

// AsImage represents packing up a directory into
// an image derived from another.
type AsImage struct {
	SourceDir         string
	TargetDir         string
	ParentImage       string
	NewRepositoryName string
	Config            docker.Config
}

func (img AsImage) WithDockerClient(c *docker.Client) error {
	imageID := utils.RandomIdentifierOfLength(64)

	var parentImageID string

	if img.ParentImage != "" {
		parentImageConfig, err := c.InspectImage(img.ParentImage)
		if err != nil {
			log.WithFields(log.Fields{"image": img.ParentImage}).Info("Parent image not found, pulling it")
			err = pull.Pull(c, img.ParentImage)
			if err != nil {
				log.WithFields(log.Fields{"image": img.ParentImage, "error": err}).Error("Pulling failed")
				return err
			}

			parentImageConfig, err = c.InspectImage(img.ParentImage)
			if err != nil {
				log.WithFields(log.Fields{"image": img.ParentImage, "error": err}).Error("Image still not found after pulling, bailing")
				return err
			}
		}
		parentImageID = parentImageConfig.ID
	} else {
		parentImageID = ""
	}

	uploadReader, uploadWriter := io.Pipe()

	layerBallName := randomTarballFileName()

	layerFile, err := os.Create(layerBallName)
	if err != nil {
		return err
	}
	//defer os.Remove(layerBallName)
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
		writeUploadBallInto(uploadWriter, layerBallName, img.NewRepositoryName, parentImageID, imageID, img.Config)
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

func writeUploadBallInto(w io.Writer, layerBallName string, newRepositoryName string, parentImageID string, imageID string, config docker.Config) error {
	uploadBall := tar.NewWriter(w)
	defer uploadBall.Close()

	repositoriesFileHeader, repoFile := repositoriesFile(newRepositoryName, imageID)
	uploadBall.WriteHeader(&repositoriesFileHeader)
	uploadBall.Write(repoFile)

	imageDirHeader := imageDir(imageID)
	uploadBall.WriteHeader(&imageDirHeader)

	versionFileHeader, versionFileContents := versionFile(imageID)
	uploadBall.WriteHeader(&versionFileHeader)
	uploadBall.Write(versionFileContents)

	configFileHeader, configFileContents := imageConfigFile(parentImageID, imageID, config)
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

func randomTarballFileName() string {
	dir := os.TempDir()
	tarid := utils.RandomIdentifier()
	return filepath.Join(dir, "involucro-volume-"+tarid+".tar")
}
