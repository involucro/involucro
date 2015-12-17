package runtask

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Runner interface {
	RunLocallyTaskWith(string, *docker.Client, string) error
	RunTaskOnRemoteSystemWith(string, *docker.Client, string) error
}

type runtaskStep struct {
	taskID string
	runner Runner
}

func (s runtaskStep) WithDockerClient(c *docker.Client, remoteWorkDir string) error {
	return s.runner.RunLocallyTaskWith(s.taskID, c, remoteWorkDir)
}

func (s runtaskStep) WithRemoteDockerClient(c *docker.Client, remoteWorkDir string) error {
	return s.runner.RunTaskOnRemoteSystemWith(s.taskID, c, remoteWorkDir)
}

func (s runtaskStep) ShowStartInfo() {
	log.WithFields(log.Fields{"ID": s.taskID}).Info("invoke task")
}
