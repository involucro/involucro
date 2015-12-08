package types

import (
	//	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
)

// Step represents one action taken by the tool.
type Step interface {
	WithDockerClient(c *docker.Client) error
}

type SubBuilder interface {
}
