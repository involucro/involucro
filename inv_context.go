package main

import "github.com/fsouza/go-dockerclient"
import duk "gopkg.in/olebedev/go-duktape.v2"

type Step interface {
	WithDockerClient(c *docker.Client) error
	DryRun()
	AsShellCommand() string
}

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
