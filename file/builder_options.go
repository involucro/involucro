package file

import (
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/thriqon/involucro/file/translator"
)

func (wbs wrapBuilderState) withConfig(l *lua.State) int {
	wbs.baseConf = translator.ParseImageConfigFromLuaTable(l)
	return wrapTable(l, &wbs)
}

func (ubs usingBuilderState) withConfig(l *lua.State) int {
	oldImageID := ubs.Config.Image
	ubs.Config = translator.ParseImageConfigFromLuaTable(l)
	if ubs.Config.Image != "" {
		log.Warn("Overwriting the used image in withConfig is discouraged")
	} else {
		ubs.Config.Image = oldImageID
	}
	return usingTable(l, &ubs)
}

func (ubs usingBuilderState) withHostConfig(l *lua.State) int {
	ubs.HostConfig = translator.ParseHostConfigFromLuaTable(l)
	return usingTable(l, &ubs)
}
