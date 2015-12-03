package file

import (
	"errors"
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
)

type RuntaskStep struct {
	context *InvContext
	taskID  string
}

func (s RuntaskStep) WithDockerClient(c *docker.Client) error {
	return s.runOrCrash(func(s Step) error {
		return s.WithDockerClient(c)
	})
}
func (s RuntaskStep) DryRun() {
	s.runOrCrash(func(s Step) error {
		s.DryRun()
		return nil
	})
}
func (s RuntaskStep) AsShellCommandOn(w io.Writer) {
	s.runOrCrash(func(s Step) error {
		s.AsShellCommandOn(w)
		return nil
	})
}

func (s RuntaskStep) runOrCrash(f func(Step) error) error {
	if s.context.HasTask(s.taskID) {
		return s.context.RunTaskWith(s.taskID, f)
	} else {
		log.WithFields(log.Fields{"task": s.taskID}).Error("Task run requested, but task not found")
		return errors.New("Task not found")
	}
}

func (inv *InvContext) HasTask(taskID string) bool {
	_, ok := inv.Tasks[taskID]
	return ok
}

func (inv *InvContext) RunTaskWith(taskID string, f func(Step) error) error {
	steps := inv.Tasks[taskID]
	if len(steps) == 0 {
		log.WithFields(log.Fields{"task": taskID}).Warn("no steps defined for task")
		return nil
	}
	for _, step := range steps {
		if err := f(step); err != nil {
			return err
		}
	}
	return nil
}

func (bs builderState) runTask(l *lua.State) int {
	otherTaskID := requireStringOrFailGracefully(l, -1, "runTask")

	rts := RuntaskStep{
		context: bs.inv,
		taskID:  otherTaskID,
	}

	tasks := bs.inv.Tasks
	tasks[bs.taskID] = append(tasks[bs.taskID], rts)
	return globalBuilderTable(l, &bs)
}
