package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

/*type Step interface {
	WithDockerClient(c *docker.Client) error
	DryRun()
	AsShellCommand() string
}*/

type ExecuteImage struct {
	opts docker.CreateContainerOptions
}

func (img ExecuteImage) DryRun() {
	log.WithFields(log.Fields{"dry": true, "image": img.opts.Config.Image}).Info("RUN")
}

func (img ExecuteImage) WithDockerClient(c *docker.Client, workdir string) error {
	_, err := c.CreateContainer(img.opts)

	if err != nil {
		log.WithFields(log.Fields{"image": img.opts.Config.Image, "err": err}).Warn("create container failed")
		return err
	}

	return nil
}

func (img ExecuteImage) AsShellCommand() string {
	return fmt.Sprintf("docker run -t --rm %s\n", img.opts.Config.Image)
}
