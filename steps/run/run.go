package run

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"regexp"
)

// ExecuteImage executes the given config and host config, similar to "docker
// run"
type ExecuteImage struct {
	Config                docker.Config
	HostConfig            docker.HostConfig
	ExpectedCode          int
	ExpectedStdoutMatcher *regexp.Regexp
	ExpectedStderrMatcher *regexp.Regexp
	ActualCode            int
}

// WithDockerClient executes the task on the given Docker instance
func (img ExecuteImage) WithDockerClient(c *docker.Client) error {

	container, err := img.createContainer(c)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"ID": container.ID}).Debug("Container created, starting it")

	if err = c.StartContainer(container.ID, nil); err != nil {
		log.WithFields(log.Fields{"ID": container.ID, "err": err}).Warn("Container not started and not removed")
		return err
	}
	log.WithFields(log.Fields{"ID": container.ID}).Debug("Container started, await completion")

	img.ActualCode, err = c.WaitContainer(container.ID)

	log.WithFields(log.Fields{"Status": img.ActualCode, "ID": container.ID}).Info("Execution complete")

	if err == nil && img.ActualCode == 0 {
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

	return err
}
