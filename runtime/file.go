package runtime

import (
	"errors"
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"os"
)

// Runtime encapsulates the state of the tool
type Runtime struct {
	lua    *lua.State
	tasks  map[string][]Step
	Values map[string]string
}

// New creates a new Runtime and returns it. This new context uses the
// working dir that is passed as a parameter.  After instantiation, the context
// will be ready to load additional files.
func New(values map[string]string) Runtime {
	m := Runtime{
		lua:    lua.NewStateEx(),
		tasks:  make(map[string][]Step),
		Values: values,
	}

	tableWith(m.lua, fm{"task": m.task})

	m.lua.SetGlobal("inv")

	m.lua.NewTable()
	tableWith(m.lua, fm{"__index": getEnvValue})
	m.lua.SetMetaTable(-2)
	m.lua.SetGlobal("ENV")

	m.lua.NewTable()
	tableWith(m.lua, fm{"__index": m.getValue})
	m.lua.SetMetaTable(-2)
	m.lua.SetGlobal("VAR")

	lua.OpenLibraries(m.lua)
	injectIoLib(m.lua)

	return m
}

// HasTask tells whether the receiver has a task with
// the given task ID, i.e. whether any steps have been
// registered for that name.
func (inv *Runtime) HasTask(taskID string) bool {
	_, ok := inv.tasks[taskID]
	return ok
}

// RunLocallyTaskWith retrieves the steps for the given task ID
// and calls f once with each step. If any error occurs
// during an invocation, this error is returned and
// the loop is interrupted.
func (inv *Runtime) RunLocallyTaskWith(taskID string, client *docker.Client, remoteWorkDir string) error {
	if !inv.HasTask(taskID) {
		log.WithFields(log.Fields{"task": taskID}).Warn("No steps defined for task")
		return errors.New("No steps defined for task")
	}

	for _, step := range inv.tasks[taskID] {
		step.ShowStartInfo()
		if err := step.WithDockerClient(client, remoteWorkDir); err != nil {
			return err
		}
	}
	return nil
}

func (inv *Runtime) RunTaskOnRemoteSystemWith(taskID string, client *docker.Client, remoteWorkDir string) error {
	if !inv.HasTask(taskID) {
		log.WithFields(log.Fields{"task": taskID}).Warn("No steps defined for task")
		return errors.New("No steps defined for task")
	}

	for _, step := range inv.tasks[taskID] {
		step.ShowStartInfo()
		if err := step.WithRemoteDockerClient(client, remoteWorkDir); err != nil {
			return err
		}
	}
	return nil
}

// RunFile runs the file with the given filename in this context
func (inv *Runtime) RunFile(fileName string) error {
	log.WithFields(log.Fields{"fileName": fileName}).Debug("Run file")
	return lua.DoFile(inv.lua, fileName)
}

// RunString runs the given parameter directly
func (inv *Runtime) RunString(script string) error {
	log.WithFields(log.Fields{"script": script}).Debug("Run script")
	return lua.DoString(inv.lua, script)
}

func getEnvValue(l *lua.State) int {
	key := lua.CheckString(l, -1)
	l.PushString(os.Getenv(key))
	return 1
}

func (inv *Runtime) getValue(l *lua.State) int {
	key := lua.CheckString(l, -1)
	l.PushString(inv.Values[key])
	return 1
}

func (inv *Runtime) task(l *lua.State) int {
	taskID := lua.CheckString(l, -1)

	registerStep := func(s Step) {
		inv.tasks[taskID] = append(inv.tasks[taskID], s)
	}

	subbuilders := make(map[string]lua.Function)
	subbuilders["task"] = inv.task

	subbuilders["using"] = newRunSubBuilder(subbuilders, registerStep)
	subbuilders["wrap"] = newWrapSubBuilder(subbuilders, registerStep)
	subbuilders["runTask"] = newRuntaskSubBuilder(subbuilders, registerStep, inv)
	subbuilders["tag"] = newTagSubBuilder(subbuilders, registerStep)

	return tableWith(l, subbuilders)
}

// TaskIDList gives back a list of tasks that are defined at the time of calling
func (inv *Runtime) TaskIDList() []string {
	taskIDs := []string{}
	for key := range inv.tasks {
		taskIDs = append(taskIDs, key)
	}
	return taskIDs
}
