package runtime

import (
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
)

type hookStep struct {
	internalHookId string
}

// ShowStartInfo displays logging information including the executed task.
func (s hookStep) ShowStartInfo() {
	log.Info("Run Hook")
}

type hookStepBuilder struct {
	hookStep
	upper        fm
	registerStep func(Step)
}

func newHookSubBuilder(upper fm, register func(Step)) lua.Function {
	hsb := hookStepBuilder{
		upper:        upper,
		registerStep: register,
	}
	return hsb.hook
}

func (hsb hookStepBuilder) hook(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsFunction(-1), 1, "expected function")

	hsb.internalHookId = randomIdentifierOfLength(20)
	l.SetField(lua.RegistryIndex, hsb.internalHookId)

	hsb.registerStep(hsb.hookStep)
	return tableWith(l, hsb.upper)
}

func (hsb hookStep) Take(i *Runtime) error {
	i.lua.Field(lua.RegistryIndex, hsb.internalHookId)
	lua.ArgumentCheck(i.lua, i.lua.IsFunction(-1), 1, "expected function as hook")

	return i.lua.ProtectedCall(0, 0, 0)
}
