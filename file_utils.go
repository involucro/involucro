
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
	if (typ == duk.TypeString) {
		return c.GetString(idx)
	} else {
		log.WithFields(log.Fields{"method": method, "type": typ}).Fatal("Invalid argument type in method call.")
		return ""
	}
}

func AddPropToObject(c *duk.Context)
