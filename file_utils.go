package main

import duk "gopkg.in/olebedev/go-duktape.v2"
import log "github.com/Sirupsen/logrus"

func PushDerivedFromThis(c *duk.Context) int {
	pos := c.PushObject()
	c.PushThis()
	c.SetPrototype(pos)
	return pos
}

func RequireStringOrFailGracefully(c *duk.Context, idx int, method string) string {
	typ := c.GetType(idx)
	if typ == duk.TypeString {
		return c.GetString(idx)
	} else {
		log.WithFields(log.Fields{"method": method, "type": typ}).Panic("Invalid argument type in method call.")
		return ""
	}
}

const DEF_PROP_FLAGS = (1 << 6) | (1 << 3)

func DefineStringOnObject(c *duk.Context, idx int, key string, value string) {
	c.PushString(key)
	c.PushString(value)
	c.DefProp(idx, DEF_PROP_FLAGS)
}

func DefineFuncOnObject(c *duk.Context, idx int, key string, value func(*duk.Context) int) {
	c.PushString(key)
	c.PushGoFunction(value)
	c.DefProp(idx, DEF_PROP_FLAGS)
}
