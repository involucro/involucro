package utils

import "github.com/fsouza/go-dockerclient"

// Step represents one action taken by the tool.
type Step interface {
	// WithDockerClient executes the receiver on the given Docker instance,
	// with non-absolute paths resolved relative to the remoteWorkDir.
	WithDockerClient(c *docker.Client, remoteWorkDir string) error

	// WithRemoteDockerClient executes the receiver on the given Docker instance,
	// assuming that the instance doesn't share the filesystem with this process.
	// Non-local paths are resolved relative to the remoteWorkDir.
	WithRemoteDockerClient(c *docker.Client, remoteWorkDir string) error

	// ShowStartInfo displays some information on the default logger that identifies the step
	ShowStartInfo()
}
