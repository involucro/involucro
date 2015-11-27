package file

import (
	"github.com/fsouza/go-dockerclient"
	duk "gopkg.in/olebedev/go-duktape.v2"
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
	duk        *duk.Context
	Tasks      map[string][]Step
	WorkingDir string
}

func (i InvContext) asCallback(f func(*InvContext) int) func(*duk.Context) int {
	return func(_ *duk.Context) int {
		return f(&i)
	}
}
