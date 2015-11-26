package main

import duk "gopkg.in/olebedev/go-duktape.v2"
import log "github.com/Sirupsen/logrus"
import "github.com/fsouza/go-dockerclient"

func InstantiateRuntimeEnv() *duk.Context {
	ctx := duk.New()

	global := ctx.PushObject()

	idx := ctx.PushObject()

	DefineFuncOnObject(ctx, idx, "task", func(c *duk.Context) int {
		task_id := RequireStringOrFailGracefully(c, -1, "task")
		log.WithFields(log.Fields{"taskId": task_id}).Info("defined task")

		retobj := PushDerivedFromThis(c)

		DefineStringOnObject(c, retobj, "taskId", task_id)
		DefineFuncOnObject(c, retobj, "using", FileUsing)
		DefineFuncOnObject(c, retobj, "wrap", FileWrapping)

		return 1
	})

	ctx.PutPropString(global, "inv")

	//idxLogger := ctx.PushObject()

	ctx.SetGlobalObject()

	return ctx
}

// WRAPPING

func FileWrapping(c *duk.Context) int {
	source_dir := RequireStringOrFailGracefully(c, -1, "wrap")
	idx := PushDerivedFromThis(c)

	DefineStringOnObject(c, idx, "sourceDir", source_dir)
	DefineFuncOnObject(c, idx, "inImage", FileInImage)
	DefineFuncOnObject(c, idx, "at", FileAt)

	return 1
}

func FileInImage(c *duk.Context) int {
	in_image := RequireStringOrFailGracefully(c, -1, "inImage")
	idx := PushDerivedFromThis(c)

	DefineStringOnObject(c, idx, "parentImage", in_image)
	return 1
}

func FileAt(c *duk.Context) int {
	target_dir := RequireStringOrFailGracefully(c, -1, "at")
	idx := PushDerivedFromThis(c)

	DefineStringOnObject(c, idx, "targetDir", target_dir)
	return 1
}

func FileAs(c *duk.Context) int {
	if c.GetTopIndex() == 2 {
		// have options parameter
	}
	return 1
}

// EXECUTION

func FileUsing(c *duk.Context) int {
	image_id := RequireStringOrFailGracefully(c, -1, "using")
	idx := PushDerivedFromThis(c)

	DefineStringOnObject(c, idx, "usingImageId", image_id)
	DefineFuncOnObject(c, idx, "run", FileRun)

	return 1
}

func FileRun(c *duk.Context) int {
	PushDerivedFromThis(c)
	c.GetPropString(-1, "taskId")
	taskId := RequireStringOrFailGracefully(c, -1, "run")
	c.Pop()

	tasks[taskId] = append(tasks[taskId], ExecuteImage{opts: docker.CreateContainerOptions{}})

	return 1
}
