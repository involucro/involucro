package utils

import "github.com/Shopify/go-lua"
import "testing"

func TestTableWith(t *testing.T) {
	l := lua.NewState()

	var a, b, c bool

	i := TableWith(l, Fm{
		"a": func(l *lua.State) int {
			a = true
			return 0
		},
	}, Fm{
		"b": func(l *lua.State) int {
			b = true
			return 0
		},
		"c": func(l *lua.State) int {
			c = true
			return 0
		},
	})
	if i != 1 {
		t.Fatal("Invalid return value from TableWith", i)
	}

	l.SetGlobal("x")

	lua.DoString(l, "x.a()")
	lua.DoString(l, "x.b()")
	lua.DoString(l, "x.c()")

	if !(a && b && c) {
		t.Fatal("Did not set all of abc to true", a, b, c)
	}
}
