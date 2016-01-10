package translator

import (
	"github.com/Shopify/go-lua"
	"testing"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if x := recover(); x == nil {
			t.Error("Did not panic")
		}
	}()
	f()
}

func TestParseDockerConfigFailures(t *testing.T) {
	t.Parallel()

	state := lua.NewState()

	run := func(s string) {
		if err := lua.DoString(state, s); err != nil {
			t.Fatal(err)
		}
	}

	assertPanic(t, func() { ParseImageConfigFromLuaTable(state) })

	run(`x = 2`)
	state.Global("x")
	assertPanic(t, func() { ParseImageConfigFromLuaTable(state) })
}
