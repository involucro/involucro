package file

import (
	"github.com/Shopify/go-lua"
	wrapS "github.com/thriqon/involucro/steps/wrap"
)

type wrapBuilderState struct {
	builderState
	wrapS.AsImage
}

func (bs builderState) wrap(l *lua.State) int {
	wbs := wrapBuilderState{
		builderState: bs,
		AsImage: wrapS.AsImage{
			SourceDir: requireStringOrFailGracefully(l, -1, "wrap"),
		},
	}

	return wrapTable(l, &wbs)
}

func (wbs wrapBuilderState) inImage(l *lua.State) int {
	wbs.ParentImage = requireStringOrFailGracefully(l, -1, "inImage")
	return wrapTable(l, &wbs)
}

func (wbs wrapBuilderState) at(l *lua.State) int {
	wbs.TargetDir = requireStringOrFailGracefully(l, -1, "at")
	return wrapTable(l, &wbs)
}

func (wbs wrapBuilderState) as(l *lua.State) int {
	wbs.NewRepositoryName = requireStringOrFailGracefully(l, -1, "as")

	tasks := wbs.inv.Tasks
	tasks[wbs.taskID] = append(tasks[wbs.taskID], wbs.AsImage)

	return wrapTable(l, &wbs)
}

func wrapTable(l *lua.State, wbs *wrapBuilderState) int {
	return tableWith(l, fm{
		"using":      wbs.using,
		"task":       wbs.task,
		"inImage":    wbs.inImage,
		"wrap":       wbs.wrap,
		"as":         wbs.as,
		"at":         wbs.at,
		"withConfig": wbs.withConfig,
	})
}
