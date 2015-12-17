package wrap

import (
	"github.com/Shopify/go-lua"
	"github.com/thriqon/involucro/file/translator"
	"github.com/thriqon/involucro/file/types"
	"github.com/thriqon/involucro/file/utils"
)

type wrapBuilderState struct {
	AsImage
	upper        utils.Fm
	registerStep func(types.Step)
}

func NewSubBuilder(upper utils.Fm, register func(types.Step)) lua.Function {
	wbs := wrapBuilderState{
		upper:        upper,
		registerStep: register,
	}
	return wbs.wrap
}

func (wbs wrapBuilderState) wrap(l *lua.State) int {
	wbs.SourceDir = lua.CheckString(l, -1)
	return wbs.wrapTable(l)
}

func (wbs wrapBuilderState) inImage(l *lua.State) int {
	wbs.ParentImage = lua.CheckString(l, -1)
	return wbs.wrapTable(l)
}

func (wbs wrapBuilderState) at(l *lua.State) int {
	wbs.TargetDir = lua.CheckString(l, -1)
	return wbs.wrapTable(l)
}

func (wbs wrapBuilderState) as(l *lua.State) int {
	wbs.NewRepositoryName = lua.CheckString(l, -1)

	wbs.registerStep(wbs.AsImage)
	return wbs.wrapTable(l)
}

func (wbs wrapBuilderState) wrapTable(l *lua.State) int {
	return utils.TableWith(l, wbs.upper, utils.Fm{
		"inImage":    wbs.inImage,
		"as":         wbs.as,
		"at":         wbs.at,
		"withConfig": wbs.withConfig,
	})
}

func (wbs wrapBuilderState) withConfig(l *lua.State) int {
	wbs.Config = translator.ParseImageConfigFromLuaTable(l)
	return wbs.wrapTable(l)
}
