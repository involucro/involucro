package runtask

import (
	"github.com/fsouza/go-dockerclient"
)

type runtaskStep struct {
	taskID        string
	runTaskWithID func(string, *docker.Client) error
}

func (s runtaskStep) WithDockerClient(c *docker.Client) error {
	return s.runTaskWithID(s.taskID, c)
}
