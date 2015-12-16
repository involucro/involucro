package types

import (
	"github.com/fsouza/go-dockerclient"
)

// Step represents one action taken by the tool.
type Step interface {
	WithDockerClient(c *docker.Client) error
	WithRemoteDockerClient(c *docker.Client) error
}
