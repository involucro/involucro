package file

import (
	"errors"
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
)

type runtaskStep struct {
	context *InvContext
	taskID  string
}

func (s runtaskStep) WithDockerClient(c *docker.Client) error {
	return s.runOrCrash(func(s Step) error {
		return s.WithDockerClient(c)
	})
}
func (s runtaskStep) DryRun() {
	s.runOrCrash(func(s Step) error {
		s.DryRun()
		return nil
	})
}
func (s runtaskStep) AsShellCommandOn(w io.Writer) {
	s.runOrCrash(func(s Step) error {
		s.AsShellCommandOn(w)
		return nil
	})
}

func (s runtaskStep) runOrCrash(f func(Step) error) error {
	if s.context.HasTask(s.taskID) {
		return s.context.RunTaskWith(s.taskID, f)
	}

	log.WithFields(log.Fields{"task": s.taskID}).Error("Task run requested, but task not found")
	return errors.New("Task not found")
}

func (bs builderState) runTask(l *lua.State) int {
	otherTaskID := requireStringOrFailGracefully(l, -1, "runTask")

	rts := runtaskStep{
		context: bs.inv,
		taskID:  otherTaskID,
	}

	tasks := bs.inv.Tasks
	tasks[bs.taskID] = append(tasks[bs.taskID], rts)
	return globalBuilderTable(l, &bs)
}
