package file

import (
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

// Step represents one action taken by the tool.
type Step interface {
	WithDockerClient(c *docker.Client) error
}

// InvContext encapsulates the state of the tool
type InvContext struct {
	lua        *lua.State
	Tasks      map[string][]Step
	WorkingDir string
}

// HasTask tells whether the receiver has a task with
// the given task ID, i.e. whether any steps have been
// registered for that name.
func (inv *InvContext) HasTask(taskID string) bool {
	_, ok := inv.Tasks[taskID]
	return ok
}

// RunTaskWith retrieves the steps for the given task ID
// and calls f once with each step. If any error occurs
// during an invocation, this error is returned and
// the loop is interrupted.
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
