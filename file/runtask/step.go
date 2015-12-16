package runtask

import (
	"github.com/fsouza/go-dockerclient"
)

type Runner interface {
	RunLocallyTaskWith(string, *docker.Client) error
	RunTaskOnRemoteSystemWith(string, *docker.Client) error
}

type runtaskStep struct {
	taskID string
	runner Runner
}

func (s runtaskStep) WithDockerClient(c *docker.Client) error {
	return s.runner.RunLocallyTaskWith(s.taskID, c)
}

func (s runtaskStep) WithRemoteDockerClient(c *docker.Client) error {
	return s.runner.RunTaskOnRemoteSystemWith(s.taskID, c)
}
