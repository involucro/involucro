package main

import duk "gopkg.in/olebedev/go-duktape.v2"
import log "github.com/Sirupsen/logrus"
import "github.com/fsouza/go-dockerclient"
import run "github.com/thriqon/involucro/steps/run"

//import "path/filepath"

import wrap "github.com/thriqon/involucro/steps/wrap"

// InstantiateRuntimeEnv creates a new InvContext and returns it. This new
// context uses the working dir that is passed as a parameter.  After
// instantiation, the context will be ready to load additional files.
func InstantiateRuntimeEnv(workingDir string) InvContext {
	m := InvContext{
		duk:        duk.New(),
		Tasks:      make(map[string][]Step),
		WorkingDir: workingDir,
	}

	global := m.duk.PushObject()

	idx := m.duk.PushObject()

	defineFuncOnObject(m.duk, idx, "task", m.asCallback(func(i *InvContext) int {
		taskID := requireStringOrFailGracefully(i.duk, -1, "task")
		log.WithFields(log.Fields{"taskId": taskID}).Info("defined task")

		retobj := pushDerivedFromThis(i.duk)

		defineStringOnObject(i.duk, retobj, "taskId", taskID)
		defineFuncOnObject(i.duk, retobj, "using", m.asCallback(fileUsing))
		defineFuncOnObject(i.duk, retobj, "wrap", m.asCallback(fileWrapping))

		return 1
	}))

	m.duk.PutPropString(global, "inv")

	//idxLogger := m.duk.PushObject()

	m.duk.SetGlobalObject()

	return m
}

// WRAPPING

func fileWrapping(c *InvContext) int {
	sourceDir := requireStringOrFailGracefully(c.duk, -1, "wrap")
	idx := pushDerivedFromThis(c.duk)

	defineStringOnObject(c.duk, idx, "sourceDir", sourceDir)
	defineFuncOnObject(c.duk, idx, "inImage", c.asCallback(fileInImage))
	defineFuncOnObject(c.duk, idx, "at", c.asCallback(fileAt))
	defineFuncOnObject(c.duk, idx, "as", c.asCallback(fileAs))

	return 1
}

func fileInImage(c *InvContext) int {
	inImage := requireStringOrFailGracefully(c.duk, -1, "inImage")
	idx := pushDerivedFromThis(c.duk)

	defineStringOnObject(c.duk, idx, "parentImage", inImage)
	return 1
}

func fileAt(c *InvContext) int {
	targetDir := requireStringOrFailGracefully(c.duk, -1, "at")
	idx := pushDerivedFromThis(c.duk)

	defineStringOnObject(c.duk, idx, "targetDir", targetDir)
	return 1
}

func fileAs(i *InvContext) int {
	c := i.duk
	newName := requireStringOrFailGracefully(c, -1, "as")
	pushDerivedFromThis(c)

	c.GetPropString(-1, "taskId")
	taskID := requireStringOrFailGracefully(c, -1, "run:task_id")
	c.Pop()

	c.GetPropString(-1, "parentImage")
	parentImage := requireStringOrFailGracefully(c, -1, "as:inImage")
	c.Pop()

	c.GetPropString(-1, "targetDir")
	targetDir := requireStringOrFailGracefully(c, -1, "as:at")
	c.Pop()

	c.GetPropString(-1, "sourceDir")
	sourceDir := requireStringOrFailGracefully(c, -1, "as:wrap")
	c.Pop()

	wi := wrap.WrapAsImage{
		SourceDir:         sourceDir,
		TargetDir:         targetDir,
		ParentImage:       parentImage,
		NewRepositoryName: newName,
	}

	i.Tasks[taskID] = append(i.Tasks[taskID], wi)

	return 1
}

// EXECUTION

func fileUsing(c *InvContext) int {
	imageID := requireStringOrFailGracefully(c.duk, -1, "using")
	idx := pushDerivedFromThis(c.duk)

	defineStringOnObject(c.duk, idx, "usingImageId", imageID)
	defineFuncOnObject(c.duk, idx, "run", c.asCallback(fileRun))

	log.WithFields(log.Fields{"usingImageId": imageID}).Debug("Using")

	return 1
}

func fileRun(i *InvContext) int {
	c := i.duk
	if c.GetTop() < 1 {
		log.WithFields(log.Fields{"method": "run"}).Panic("Required argument missing")
	}

	topIndex := c.GetTopIndex()

	cmd := make([]string, topIndex+1)
	for pos := 0; pos <= topIndex; pos++ {
		cmd[pos] = requireStringOrFailGracefully(c, pos, "run")
	}

	pushDerivedFromThis(c)
	c.GetPropString(-1, "taskId")
	taskID := requireStringOrFailGracefully(c, -1, "run:task_id")
	c.Pop()
	c.GetPropString(-1, "usingImageId")
	imageID := requireStringOrFailGracefully(c, -1, "run:image_id")
	c.Pop()

	bind := i.WorkingDir + ":/source"

	ei := run.ExecuteImage{
		Config: docker.Config{
			Cmd:   cmd,
			Image: imageID,
		},
		HostConfig: docker.HostConfig{
			Binds: []string{
				bind,
			},
		},
	}

	i.Tasks[taskID] = append(i.Tasks[taskID], ei)
	return 1
}
