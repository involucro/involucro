package file

import (
	"github.com/Shopify/go-lua"
)

// InstantiateRuntimeEnv creates a new InvContext and returns it. This new
// context uses the working dir that is passed as a parameter.  After
// instantiation, the context will be ready to load additional files.
func InstantiateRuntimeEnv(workingDir string) InvContext {
	m := InvContext{
		lua:        lua.NewStateEx(),
		Tasks:      make(map[string][]Step),
		WorkingDir: workingDir,
	}

	tableWith(m.lua, fm{"task": m.task})
	m.lua.SetGlobal("inv")

	return m
}
