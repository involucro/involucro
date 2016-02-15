// **Runtask** implements the step protocol by invoking another task defined in
// the same file.  This can be used to split up the steps to the final result
// into multiple and smaller tasks that can be invoked separataly, but also
// providing an overall task comprising them.

package runtime

import (
	"github.com/Shopify/go-lua"
	"github.com/thriqon/involucro/ilog"
)

// ## The Step
// The step structure just contains the ID of the task that is executed as part
// of the step, and an interface to execute that task. The interface is defined
// below.
type runtaskStep struct {
	taskID string
}

func (s runtaskStep) Take(i *Runtime) error {
	return i.RunTask(s.taskID)
}

// ShowStartInfo displays logging information including the executed task.
func (s runtaskStep) ShowStartInfo() {
	ilog.Info.Logf("Invoke Task [%s]", s.taskID)
}

// ## The Builder
// The interface to the builder is very simple. There is only one method,
// which is both introductory and final method at the same time: "runTask".
// As such, the required builder state only stores the resulting step,
// the upper Lua function table and a method to register steps.
type runtaskBuilderState struct {
	runtaskStep
	upper        fm
	registerStep func(Step)
}

func newRuntaskSubBuilder(upper fm, register func(Step)) lua.Function {
	rbs := runtaskBuilderState{
		upper:        upper,
		registerStep: register,
	}
	return rbs.runTask
}

func (rbs runtaskBuilderState) runTask(l *lua.State) int {
	otherTaskID := lua.CheckString(l, -1)

	rbs.taskID = otherTaskID

	rbs.registerStep(rbs.runtaskStep)
	return tableWith(l, rbs.upper)
}
