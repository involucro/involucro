package run

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/file/utils"
)

func (img ExecuteImage) createContainer(c *docker.Client) (container *docker.Container, err error) {
	containerName := "step-" + utils.RandomIdentifier()

	opts := docker.CreateContainerOptions{
		Name:       containerName,
		Config:     &img.Config,
		HostConfig: &img.HostConfig,
	}

	log.WithFields(log.Fields{"containerName": containerName}).Debug("Create Container")
	container, err = c.CreateContainer(opts)

	if err == docker.ErrNoSuchImage {
		if err = utils.Pull(c, img.Config.Image); err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("pull failed")
			return
		}

		log.WithFields(log.Fields{"containerName": containerName}).Debug("Retry: Create Container")
		container, err = c.CreateContainer(opts)
	}

	if err != nil {
		log.WithFields(log.Fields{"image": img.Config.Image, "err": err}).Warn("create container failed")
	}
	return
}
