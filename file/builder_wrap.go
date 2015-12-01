package file

import (
	"github.com/Shopify/go-lua"
	wrapS "github.com/thriqon/involucro/steps/wrap"
)

type wrapBuilderState struct {
	builderState
	sourceDir     string
	targetDir     string
	parentImageID string
}

func (bs builderState) wrap(l *lua.State) int {
	wbs := wrapBuilderState{
		builderState: bs,
		sourceDir:    requireStringOrFailGracefully(l, -1, "wrap"),
	}

	return wrapTable(l, &wbs)
}

func (wbs wrapBuilderState) inImage(l *lua.State) int {
	wbsn := wbs
	wbsn.parentImageID = requireStringOrFailGracefully(l, -1, "inImage")

	return wrapTable(l, &wbsn)
}

func (wbs wrapBuilderState) at(l *lua.State) int {
	wbsn := wbs
	wbsn.targetDir = requireStringOrFailGracefully(l, -1, "at")
	return wrapTable(l, &wbsn)
}

func (wbs wrapBuilderState) as(l *lua.State) int {
	ai := wrapS.AsImage{
		SourceDir:         wbs.sourceDir,
		TargetDir:         wbs.targetDir,
		ParentImage:       wbs.parentImageID,
		NewRepositoryName: requireStringOrFailGracefully(l, -1, "as"),
	}

	tasks := wbs.inv.Tasks
	tasks[wbs.taskID] = append(tasks[wbs.taskID], ai)

	return wrapTable(l, &wbs)
}

func wrapTable(l *lua.State, wbs *wrapBuilderState) int {
	return tableWith(l, fm{
		"using":   wbs.using,
		"task":    wbs.task,
		"inImage": wbs.inImage,
		"wrap":    wbs.wrap,
		"as":      wbs.as,
		"at":      wbs.at,
	})
}
