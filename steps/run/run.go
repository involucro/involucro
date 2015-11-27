package run

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	utils "github.com/thriqon/involucro/lib"
	"io"
)

// Implementation of the Step interface
// Executes the given config and host config, similar to "docker run"
type ExecuteImage struct {
	Config     docker.Config
	HostConfig docker.HostConfig
}

// DryRun runs this task without doing anything, but logging an indication of
// what would have been done
func (img ExecuteImage) DryRun() {
	log.WithFields(log.Fields{"dry": true, "image": img.Config.Image}).Info("RUN")
}

// WithDockerClient executes the task on the given Docker instance
func (img ExecuteImage) WithDockerClient(c *docker.Client) error {
	var err error
	var container *docker.Container

	containerName := "step-" + utils.RandomIdentifier()

	opts := docker.CreateContainerOptions{
		Name:       containerName,
		Config:     &img.Config,
		HostConfig: &img.HostConfig,
	}
	log.WithFields(log.Fields{"containerName": containerName}).Debug("Create Container")
	container, err = c.CreateContainer(opts)

	if err == docker.ErrNoSuchImage {
		err = pull(c, img.Config.Image)

		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("pull failed")
			return err
		}

		log.WithFields(log.Fields{"containerName": containerName}).Debug("Retry: Create Container")
		container, err = c.CreateContainer(opts)
	}

	if err != nil {
		log.WithFields(log.Fields{"image": img.Config.Image, "err": err}).Warn("create container failed")
		return err
	}

	log.WithFields(log.Fields{"ID": container.ID}).Debug("Container created, starting it")
	err = c.StartContainer(container.ID, nil)

	if err != nil {
		log.WithFields(log.Fields{"ID": container.ID, "err": err}).Warn("Container not started and not removed")
		return err
	} else {
		log.WithFields(log.Fields{"ID": container.ID}).Debug("Container started, await completion")
	}

	status, waitErr := c.WaitContainer(container.ID)

	log.WithFields(log.Fields{"Status": status, "ID": container.ID}).Info("Execution complete")

	if waitErr == nil && status == 0 {
		err := c.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
		if err != nil {
			log.WithFields(log.Fields{"ID": container.ID, "err": err}).Warn("Container not removed")
		} else {
			log.WithFields(log.Fields{"ID": container.ID}).Debug("Container removed")
		}
	} else {
		log.Debug("There was an error in execution or creation, container not removed")
	}

	return waitErr
}

// AsShellCommandOn prints sh compatible commands into the given writer, that
// accomplish the funciontality encoded in this step
func (img ExecuteImage) AsShellCommandOn(w io.Writer) {
	fmt.Fprintf(w, "docker run -t --rm %s\n", img.Config.Image)
}
