package runtask

import (
	"github.com/Shopify/go-lua"
	"github.com/thriqon/involucro/file/utils"
)

type runtaskBuilderState struct {
	runtaskStep
	upper        utils.Fm
	registerStep func(utils.Step)
}

func NewSubBuilder(upper utils.Fm, register func(utils.Step), runner Runner) lua.Function {
	rbs := runtaskBuilderState{
		runtaskStep: runtaskStep{
			runner: runner,
		},
		upper:        upper,
		registerStep: register,
	}
	return rbs.runTask
}

func (rbs runtaskBuilderState) runTask(l *lua.State) int {
	otherTaskID := lua.CheckString(l, -1)

	rbs.taskID = otherTaskID

	rbs.registerStep(rbs.runtaskStep)
	return utils.TableWith(l, rbs.upper)
}
