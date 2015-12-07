package file

import (
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
)

func requireStringOrFailGracefully(c *lua.State, idx int, method string) string {
	if !c.IsString(idx) {
		log.WithFields(log.Fields{"method": method}).Panic("Invalid argument type in method call.")
		return ""
	}
	str, _ := c.ToString(idx)
	return str
}

type fm map[string]lua.Function

func tableWith(l *lua.State, fs fm) int {
	l.CreateTable(0, len(fs))
	idx := l.Top()
	for k := range fs {
		l.PushGoFunction(fs[k])
		l.SetField(idx, k)
	}
	return 1
}

func argumentsToStringArray(l *lua.State) (args []string) {
	top := l.Top()
	args = make([]string, top)
	for i := 1; i <= top; i++ {
		args[i-1] = requireStringOrFailGracefully(l, i, "run")
	}
	return
}
