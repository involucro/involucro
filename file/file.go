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
	lua    *lua.State
	Tasks  map[string][]types.Step
	Values map[string]string
}

// InstantiateRuntimeEnv creates a new InvContext and returns it. This new
// context uses the working dir that is passed as a parameter.  After
// instantiation, the context will be ready to load additional files.
func InstantiateRuntimeEnv(values map[string]string) InvContext {
	m := InvContext{
		lua:    lua.NewStateEx(),
		Tasks:  make(map[string][]types.Step),
		Values: values,
	}

	utils.TableWith(m.lua, utils.Fm{"task": m.task})

	m.lua.SetGlobal("inv")

	m.lua.NewTable()
	utils.TableWith(m.lua, utils.Fm{"__index": getEnvValue})
	m.lua.SetMetaTable(-2)
	m.lua.SetGlobal("ENV")

	m.lua.NewTable()
	utils.TableWith(m.lua, utils.Fm{"__index": m.getValue})
	m.lua.SetMetaTable(-2)
	m.lua.SetGlobal("VAR")

	return m
}

// HasTask tells whether the receiver has a task with
// the given task ID, i.e. whether any steps have been
// registered for that name.
func (inv *InvContext) HasTask(taskID string) bool {
	_, ok := inv.Tasks[taskID]
	return ok
}

// RunLocallyTaskWith retrieves the steps for the given task ID
// and calls f once with each step. If any error occurs
// during an invocation, this error is returned and
// the loop is interrupted.
func (inv *InvContext) RunLocallyTaskWith(taskID string, client *docker.Client, remoteWorkDir string) error {
	if !inv.HasTask(taskID) {
		log.WithFields(log.Fields{"task": taskID}).Warn("No steps defined for task")
		return errors.New("No steps defined for task")
	}

	for _, step := range inv.Tasks[taskID] {
		if err := step.WithDockerClient(client, remoteWorkDir); err != nil {
			return err
		}
	}
	return nil
}

func (inv *InvContext) RunTaskOnRemoteSystemWith(taskID string, client *docker.Client, remoteWorkDir string) error {
	if !inv.HasTask(taskID) {
		log.WithFields(log.Fields{"task": taskID}).Warn("No steps defined for task")
		return errors.New("No steps defined for task")
	}

	for _, step := range inv.Tasks[taskID] {
		if err := step.WithRemoteDockerClient(client, remoteWorkDir); err != nil {
			return err
		}
	}
	return nil
}

// RunFile runs the file with the given filename in this context
func (inv *InvContext) RunFile(fileName string) error {
	log.WithFields(log.Fields{"fileName": fileName}).Debug("Run file")
	return lua.DoFile(inv.lua, fileName)
}

// RunString runs the given parameter directly
func (inv *InvContext) RunString(script string) error {
	log.WithFields(log.Fields{"script": script}).Debug("Run script")
	return lua.DoString(inv.lua, script)
}

func getEnvValue(l *lua.State) int {
	key := lua.CheckString(l, -1)
	l.PushString(os.Getenv(key))
	return 1
}

func (inv *InvContext) getValue(l *lua.State) int {
	key := lua.CheckString(l, -1)
	l.PushString(inv.Values[key])
	return 1
}
