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

type WrapAsImage struct {
}

func (img WrapAsImage) DryRun() {
	log.WithFields(log.Fields{"dry": true}).Info("WRAP")
}

func (img WrapAsImage) WithDockerClient(c *docker.Client) error {
	return nil
}

func (img WrapAsImage) AsShellCommand() string {
	return fmt.Sprintf("docker wrap -t --rm \n")
}
