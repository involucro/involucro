package wrap

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/file/run"
	"github.com/thriqon/involucro/file/types"
	"github.com/thriqon/involucro/file/utils"
)

func init() {
	utils.RegisterEncodeableType(AsImage{})
}

func (ai AsImage) forRemoteExecution() types.Step {
	dockerSocket := "/var/run/docker.sock"

	origSourceDir := ai.SourceDir
	ai.SourceDir = "/source"

	steps := []types.Step{ai}

	return run.ExecuteImage{
		Config: docker.Config{
			Image: "involucro/tool:latest",
			Env:   []string{"STATE=" + utils.EncodeState(steps)},
			Cmd:   []string{"--encoded-state", "--socket", "/sock"},
		},
		HostConfig: docker.HostConfig{
			Binds: []string{
				dockerSocket + ":/sock",
				origSourceDir + ":/source",
			},
		},
	}
}

func (ai AsImage) WithRemoteDockerClient(c *docker.Client) error {
	return ai.forRemoteExecution().WithRemoteDockerClient(c)
}
