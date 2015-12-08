package utils

import "github.com/Shopify/go-lua"

type Fm map[string]lua.Function

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
