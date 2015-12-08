package translator

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
)

func checkBoolean(l *lua.State, index int) bool {
	lua.CheckType(l, index, lua.TypeBoolean)
	return l.ToBoolean(index)
}

func checkStringArray(l *lua.State, index int) []string {
	lua.CheckType(l, index, lua.TypeTable)

	var items []string

	l.PushNil()
	for l.Next(-2) {
		items = append(items, lua.CheckString(l, -1))
		l.Pop(1)
	}

	return items
}

func checkStringMap(l *lua.State, index int) map[string]string {
	lua.CheckType(l, index, lua.TypeTable)
	items := make(map[string]string)

	l.PushNil()
	for l.Next(-2) {
		key := lua.CheckString(l, -2)
		val := lua.CheckString(l, -1)
		items[key] = val
		l.Pop(1)
	}

	return items
}

func parseExposedPorts(l *lua.State, index int) map[docker.Port]struct{} {
	items := map[docker.Port]struct{}{}

	for _, el := range checkStringArray(l, index) {
		items[docker.Port(el)] = struct{}{}
	}

	return items
}

func checkStringSet(l *lua.State, index int) map[string]struct{} {
	items := map[string]struct{}{}
	for _, el := range checkStringArray(l, index) {
		items[el] = struct{}{}
	}
	return items
}
