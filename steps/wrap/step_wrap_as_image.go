package wrap

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type WrapAsImage struct {
	SourceDir         string
	TargetDir         string
	ParentImage       string
	NewRepositoryName string
}

func (img WrapAsImage) DryRun() {
	log.WithFields(log.Fields{"dry": true}).Info("WRAP")
}

func (img WrapAsImage) WithDockerClient(c *docker.Client) error {
	return nil
}

func (img WrapAsImage) AsShellCommand() string {
	return fmt.Sprintf("echo NOT IMPLEMENTED YET\n")
}
