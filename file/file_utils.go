package file

import duk "gopkg.in/olebedev/go-duktape.v2"
import log "github.com/Sirupsen/logrus"

func pushDerivedFromThis(c *duk.Context) int {
	pos := c.PushObject()
	c.PushThis()
	c.SetPrototype(pos)
	return pos
}

func requireStringOrFailGracefully(c *duk.Context, idx int, method string) string {
	typ := c.GetType(idx)
	if typ != duk.TypeString {
		log.WithFields(log.Fields{"method": method, "type": typ}).Panic("Invalid argument type in method call.")
		return ""
	}
	return c.GetString(idx)

}

const defPropFlags = (1 << 6) | (1 << 3)

func defineStringOnObject(c *duk.Context, idx int, key string, value string) {
	c.PushString(key)
	c.PushString(value)
	c.DefProp(idx, defPropFlags)
}

func defineFuncOnObject(c *duk.Context, idx int, key string, value func(*duk.Context) int) {
	c.PushString(key)
	c.PushGoFunction(value)
	c.DefProp(idx, defPropFlags)
}
