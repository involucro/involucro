package main

import duk "gopkg.in/olebedev/go-duktape.v2"
import log "github.com/Sirupsen/logrus"
import "github.com/fsouza/go-dockerclient"
import run "github.com/thriqon/involucro/steps/run"

//import wrap "github.com/thriqon/involucro/steps/wrap"

func InstantiateRuntimeEnv(workingDir string) InvContext {
	m := InvContext{
		duk:        duk.New(),
		Tasks:      make(map[string][]Step),
		WorkingDir: workingDir,
	}

	global := m.duk.PushObject()

	idx := m.duk.PushObject()

	DefineFuncOnObject(m.duk, idx, "task", m.asCallback(func(i *InvContext) int {
		task_id := RequireStringOrFailGracefully(i.duk, -1, "task")
		log.WithFields(log.Fields{"taskId": task_id}).Info("defined task")

		retobj := PushDerivedFromThis(i.duk)

		DefineStringOnObject(i.duk, retobj, "taskId", task_id)
		DefineFuncOnObject(i.duk, retobj, "using", m.asCallback(FileUsing))
		DefineFuncOnObject(i.duk, retobj, "wrap", m.asCallback(FileWrapping))

		return 1
	}))

	m.duk.PutPropString(global, "inv")

	//idxLogger := m.duk.PushObject()

	m.duk.SetGlobalObject()

	return m
}

// WRAPPING

func FileWrapping(c *InvContext) int {
	source_dir := RequireStringOrFailGracefully(c.duk, -1, "wrap")
	idx := PushDerivedFromThis(c.duk)

	DefineStringOnObject(c.duk, idx, "sourceDir", source_dir)
	DefineFuncOnObject(c.duk, idx, "inImage", c.asCallback(FileInImage))
	DefineFuncOnObject(c.duk, idx, "at", c.asCallback(FileAt))

	return 1
}

func FileInImage(c *InvContext) int {
	in_image := RequireStringOrFailGracefully(c.duk, -1, "inImage")
	idx := PushDerivedFromThis(c.duk)

	DefineStringOnObject(c.duk, idx, "parentImage", in_image)
	return 1
}

func FileAt(c *InvContext) int {
	target_dir := RequireStringOrFailGracefully(c.duk, -1, "at")
	idx := PushDerivedFromThis(c.duk)

	DefineStringOnObject(c.duk, idx, "targetDir", target_dir)
	return 1
}

func FileAs(c *InvContext) int {
	if c.duk.GetTopIndex() == 2 {
		// have options parameter
	}
	return 1
}

// EXECUTION

func FileUsing(c *InvContext) int {
	image_id := RequireStringOrFailGracefully(c.duk, -1, "using")
	idx := PushDerivedFromThis(c.duk)

	DefineStringOnObject(c.duk, idx, "usingImageId", image_id)
	DefineFuncOnObject(c.duk, idx, "run", c.asCallback(FileRun))

	log.WithFields(log.Fields{"usingImageId": image_id}).Debug("Using")

	return 1
}

func FileRun(i *InvContext) int {
	c := i.duk
	if c.GetTop() < 1 {
		log.WithFields(log.Fields{"method": "run"}).Panic("Required argument missing")
	}

	top_index := c.GetTopIndex()

	cmd := make([]string, top_index+1)
	for pos := 0; pos <= top_index; pos++ {
		cmd[pos] = RequireStringOrFailGracefully(c, pos, "run")
	}

	PushDerivedFromThis(c)
	c.GetPropString(-1, "taskId")
	taskId := RequireStringOrFailGracefully(c, -1, "run:task_id")
	c.Pop()
	c.GetPropString(-1, "usingImageId")
	imageId := RequireStringOrFailGracefully(c, -1, "run:image_id")
	c.Pop()

	bind := i.WorkingDir + ":/source"

	ei := run.ExecuteImage{
		Config: docker.Config{
			Cmd:   cmd,
			Image: imageId,
		},
		HostConfig: docker.HostConfig{
			Binds: []string{
				bind,
			},
		},
	}

	i.Tasks[taskId] = append(i.Tasks[taskId], ei)
	return 1
}
