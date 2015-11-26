
package main

import duk "gopkg.in/olebedev/go-duktape.v2"
import log "github.com/Sirupsen/logrus"

func InstantiateRuntimeEnv() *duk.Context {
	ctx := duk.New()

	global := ctx.PushObject()

	idx := ctx.PushObject()
	ctx.PushGoFunction(func (c *duk.Context) int {
		task_id := RequireStringOrFailGracefully(c, -1, "task")
		log.WithFields(log.Fields{"taskId": task_id}).Info("defined task")
		retval := PushDerivedFromThis(c)

		c.PushString(task_id)
		c.PutPropString(retval, "taskId")

		return 1
	})
	ctx.PutPropString(idx, "task")

	ctx.PutPropString(global, "inv")

	//idxLogger := ctx.PushObject()

	ctx.SetGlobalObject()


	return ctx
}

func FileUsing(c *duk.Context) int {
	image_id := RequireStringOrFailGracefully(c, -1, "using")
	idx := PushDerivedFromThis(c)

	c.PushString(image_id)
	c.PutPropString(idx, "usingImageId")

	return 1
}

