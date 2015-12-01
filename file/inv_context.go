package file

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"io"
)

// Step represents one action taken by the tool.
type Step interface {
	WithDockerClient(c *docker.Client) error
	DryRun()
	AsShellCommandOn(w io.Writer)
}

// InvContext encapsulates the state of the tool
type InvContext struct {
	lua        *lua.State
	Tasks      map[string][]Step
	WorkingDir string
}
