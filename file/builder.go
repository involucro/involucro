package file

import (
	"github.com/Shopify/go-lua"

	"github.com/thriqon/involucro/file/run"
	"github.com/thriqon/involucro/file/runtask"
	"github.com/thriqon/involucro/file/wrap"

	"github.com/thriqon/involucro/file/types"
	"github.com/thriqon/involucro/file/utils"
)

func (inv InvContext) task(l *lua.State) int {
	taskID := lua.CheckString(l, -1)

	registerStep := func(s types.Step) {
		inv.Tasks[taskID] = append(inv.Tasks[taskID], s)
	}

	subbuilders := make(map[string]lua.Function)
	subbuilders["task"] = inv.task

	subbuilders["using"] = run.NewSubBuilder(subbuilders, registerStep, inv.WorkingDir)
	subbuilders["wrap"] = wrap.NewSubBuilder(subbuilders, registerStep, inv.WorkingDir)
	subbuilders["runTask"] = runtask.NewSubBuilder(subbuilders, registerStep, inv.RunTaskWith)

	return utils.TableWith(l, subbuilders)
}
