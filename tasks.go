package main

import "github.com/fsouza/go-dockerclient"

type Step interface {
	WithDockerClient(c *docker.Client, workdir string) error
	DryRun()
	AsShellCommand() string
}

var tasks = make(map[string][]Step)
