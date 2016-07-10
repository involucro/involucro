package runtime

import (
	"os"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"github.com/involucro/involucro/ilog"
)

// Runtime encapsulates the state of the tool
type Runtime struct {
	lua     *lua.State
	tasks   map[string][]Step
	Values  map[string]string
	client  *docker.Client
	workDir string
}

// New creates a new Runtime and returns it. This new context uses the
// working dir that is passed as a parameter.  After instantiation, the context
// will be ready to load additional files.
func New(values map[string]string, c *docker.Client, workDir string) Runtime {
	m := Runtime{
		lua:     lua.NewStateEx(),
		tasks:   make(map[string][]Step),
		Values:  values,
		client:  c,
		workDir: workDir,
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

func (inv *Runtime) isUsingRemoteInstance() bool {
	return strings.HasPrefix(inv.client.Endpoint(), "tcp:")
}

// HasTask tells whether the receiver has a task with
// the given task ID, i.e. whether any steps have been
// registered for that name.
func (inv *Runtime) HasTask(taskID string) bool {
	_, ok := inv.tasks[taskID]
	return ok
}

// RunTask executes the task with the given ID.  It is perfectly fine to
// execute a task without any steps, but a warning will be shown. Any error
// that occurs during task processing is returned.
func (inv *Runtime) RunTask(taskID string) error {
	steps := inv.tasks[taskID]
	if len(steps) == 0 {
		ilog.Warn.Logf("Invoked Stepless Task [%s]", taskID)
		return nil
	}
	for _, step := range steps {
		step.ShowStartInfo()
		if err := step.Take(inv); err != nil {
			return err
		}
	}
	return nil
}

// RunFile runs the file with the given filename in this context
func (inv *Runtime) RunFile(fileName string) error {
	ilog.Debug.Logf("Run file [%s]", fileName)
	return lua.DoFile(inv.lua, fileName)
}

// RunString runs the given parameter directly
func (inv *Runtime) RunString(script string) error {
	ilog.Debug.Logf("Run script [%s]", script)
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
	subbuilders["runTask"] = newRuntaskSubBuilder(subbuilders, registerStep)
	subbuilders["tag"] = newTagSubBuilder(subbuilders, registerStep)
	subbuilders["hook"] = newHookSubBuilder(subbuilders, registerStep)
	subbuilders["push"] = newPushSubBuilder(subbuilders, registerStep)

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
