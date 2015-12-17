package utils

import "github.com/Shopify/go-lua"

// Fm is used as a blueprint for a Lua table,
// containing functions identified by strings.
type Fm map[string]lua.Function

// TableWith translates the given blueprint to
// a table, pushed to the top of the stack.
func TableWith(l *lua.State, fs ...Fm) int {
	l.CreateTable(0, len(fs))
	idx := l.Top()
	for _, x := range fs {
		for k := range x {
			l.PushGoFunction(x[k])
			l.SetField(idx, k)
		}
	}
	return 1
}
