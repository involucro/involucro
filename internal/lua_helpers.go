package runtime

import "github.com/Shopify/go-lua"

// fm is used as a blueprint for a Lua table,
// containing functions identified by strings.
type fm map[string]lua.Function

// tableWith translates the given blueprint to
// a table, pushed to the top of the stack.
func tableWith(l *lua.State, fs ...fm) int {
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
