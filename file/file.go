package file

import (
	"errors"
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/file/types"
	"github.com/thriqon/involucro/file/utils"
	"os"
)

// InvContext encapsulates the state of the tool
type InvContext struct {
	lua        *lua.State
	Tasks      map[string][]types.Step
	WorkingDir string
}

// InstantiateRuntimeEnv creates a new InvContext and returns it. This new
// context uses the working dir that is passed as a parameter.  After
// instantiation, the context will be ready to load additional files.
func InstantiateRuntimeEnv(workingDir string) InvContext {
	m := InvContext{
		lua:        lua.NewStateEx(),
		Tasks:      make(map[string][]types.Step),
		WorkingDir: workingDir,
	}

	utils.TableWith(m.lua, utils.Fm{"task": m.task})

	m.lua.SetGlobal("inv")

	m.lua.NewTable()
	utils.TableWith(m.lua, utils.Fm{"__index": getEnvValue})
	m.lua.SetMetaTable(-2)
	m.lua.SetGlobal("ENV")

	return m
}

// HasTask tells whether the receiver has a task with
// the given task ID, i.e. whether any steps have been
// registered for that name.
func (inv *InvContext) HasTask(taskID string) bool {
	_, ok := inv.Tasks[taskID]
	return ok
}

// RunTaskWith retrieves the steps for the given task ID
// and calls f once with each step. If any error occurs
// during an invocation, this error is returned and
// the loop is interrupted.
func (inv *InvContext) RunTaskWith(taskID string, client *docker.Client) error {
	steps := inv.Tasks[taskID]
	if len(steps) == 0 {
		log.WithFields(log.Fields{"task": taskID}).Warn("no steps defined for task")
		return errors.New("No steps defined for task")
	}
	for _, step := range steps {
		if err := step.WithDockerClient(client); err != nil {
			return err
		}
	}
	return nil
}

// RunFile runs the file with the given filename in this context
func (i *InvContext) RunFile(fileName string) error {
	log.WithFields(log.Fields{"fileName": fileName}).Debug("Run file")
	return lua.DoFile(i.lua, fileName)
}

// RunString runs the given parameter directly
func (i *InvContext) RunString(script string) error {
	log.WithFields(log.Fields{"script": script}).Debug("Run script")
	return lua.DoString(i.lua, script)
}

func getEnvValue(l *lua.State) int {
	key := lua.CheckString(l, -1)
	l.PushString(os.Getenv(key))
	return 1
}
