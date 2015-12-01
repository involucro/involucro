package file

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	runS "github.com/thriqon/involucro/steps/run"
)

type usingBuilderState struct {
	builderState
	imageID string
}

func (bs builderState) using(l *lua.State) int {
	nbs := usingBuilderState{
		builderState: bs,
		imageID:      requireStringOrFailGracefully(l, -1, "using"),
	}
	return usingTable(l, &nbs)
}

func (ubs usingBuilderState) usingRun(l *lua.State) int {
	args := argumentsToStringArray(l)
	ei := runS.ExecuteImage{
		Config: docker.Config{
			Image:      ubs.imageID,
			Cmd:        args,
			WorkingDir: "/source",
		},
		HostConfig: docker.HostConfig{
			Binds: []string{
				ubs.inv.WorkingDir + ":/source",
			},
		},
	}
	tasks := ubs.inv.Tasks
	tasks[ubs.taskID] = append(tasks[ubs.taskID], ei)

	return usingTable(l, &ubs)
}

func usingTable(l *lua.State, ubs *usingBuilderState) int {
	return tableWith(l, fm{
		"using": ubs.using,
		"run":   ubs.usingRun,
		"task":  ubs.task,
	})
}
