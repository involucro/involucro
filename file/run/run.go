package run

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"os"
	"path"
	"regexp"
	"strings"
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

func (img ExecuteImage) WithRemoteDockerClient(c *docker.Client, remoteWorkDir string) error {
	if !path.IsAbs(remoteWorkDir) {
		remoteWorkDir = path.Join("/", remoteWorkDir)
	}
	return img.withAbsolutizedWorkDir(c, remoteWorkDir)
}

func (img ExecuteImage) WithDockerClient(c *docker.Client, remoteWorkDir string) error {
	if !path.IsAbs(remoteWorkDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		remoteWorkDir = path.Join(cwd, remoteWorkDir)
	}
	return img.withAbsolutizedWorkDir(c, remoteWorkDir)
}

func (img ExecuteImage) withAbsolutizedWorkDir(c *docker.Client, remoteWorkDir string) error {
	img.HostConfig = absolutizeBinds(img.HostConfig, remoteWorkDir)

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

	err = img.loadAndProcessLogs(c, container.ID)
	if err != nil {
		return err
	}

	img.ActualCode, err = c.WaitContainer(container.ID)

	if img.ActualCode != img.ExpectedCode {
		log.WithFields(log.Fields{"ID": container.ID, "expected": img.ExpectedCode, "actual": img.ActualCode}).Error("Unexpected exit code, container not removed")
		return errors.New("Unexpected exit code")
	}

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

func absolutizeBinds(h docker.HostConfig, workDir string) docker.HostConfig {
	for ind, el := range h.Binds {
		parts := strings.Split(el, ":")
		if len(parts) != 2 {
			log.WithFields(log.Fields{"bind": el}).Panic("Invalid bind, has to be of the form: source:dest")
		}

		if !path.IsAbs(parts[0]) {
			h.Binds[ind] = path.Join(workDir, parts[0]) + ":" + parts[1]
		}
	}
	return h
}
