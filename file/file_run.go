package file

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	run "github.com/thriqon/involucro/steps/run"
)

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
			Cmd:        cmd,
			Image:      imageID,
			WorkingDir: "/source",
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
