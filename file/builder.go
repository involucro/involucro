package file

import (
	"github.com/Shopify/go-lua"
)

type builderState struct {
	task   lua.Function
	scope  string
	taskID string
	inv    *InvContext
}

func (inv *InvContext) task(l *lua.State) int {
	bs := builderState{
		taskID: requireStringOrFailGracefully(l, -1, "task"),
		inv:    inv,
	}

	return globalBuilderTable(l, &bs)
}

func globalBuilderTable(l *lua.State, bs *builderState) int {
	return tableWith(l, fm{
		"using":   bs.using,
		"wrap":    bs.wrap,
		"runTask": bs.runTask,
	})
}
