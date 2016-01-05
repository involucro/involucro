// **Runtask** implements the step protocol by invoking another task defined in
// the same file.  This can be used to split up the steps to the final result
// into multiple and smaller tasks that can be invoked separataly, but also
// providing an overall task comprising them.

package runtime

import (
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

// ## The Step
// The step structure just contains the ID of the task that is executed as part
// of the step, and an interface to execute that task. The interface is defined
// below.
type runtaskStep struct {
	taskID string
	runner taskRunner
}

// The task runner interface provides the possibility to run tasks on a
// Docker client either locally or remotely.
type taskRunner interface {
	// RunLocallyTaskWith executes the task with the id of the first parameter on
	// the client given by the second parameter. The remote working dir is given
	// in the third parameter. Any error during that call is returned, otherwise
	// it is nil.  This method assumes that the Docker instance is running in the
	// same filesystem and networking stack as the calling process. In contrast,
	// the variant RunTaskOnRemoteSystemWith does not make this assumption.
	RunLocallyTaskWith(string, *docker.Client, string) error

	// RunTaskOnRemoteSystemWith is similar to RunLocallyTaskWith, see there.
	RunTaskOnRemoteSystemWith(string, *docker.Client, string) error
}

// WithDockerClient executes the task with the ID stored in the step on the
// given Docker client.
func (s runtaskStep) WithDockerClient(c *docker.Client, remoteWorkDir string) error {
	return s.runner.RunLocallyTaskWith(s.taskID, c, remoteWorkDir)
}

// WithRemoteDockerClient executes the task with the ID stored in the step on
// the given Docker client, assuming it resides in a remote filesystem and
// networking space.
func (s runtaskStep) WithRemoteDockerClient(c *docker.Client, remoteWorkDir string) error {
	return s.runner.RunTaskOnRemoteSystemWith(s.taskID, c, remoteWorkDir)
}

// ShowStartInfo displays logging information including the executed task.
func (s runtaskStep) ShowStartInfo() {
	log.WithFields(log.Fields{"ID": s.taskID}).Info("Invoke Task")
}

// ## The Builder
// The interface to the builder is very simple. There is only one method,
// which is both introductory and final method at the same time: "runTask".
// As such, the required builder state only stores the resulting step,
// the upper Lua function table and a method to register steps.
type runtaskBuilderState struct {
	runtaskStep
	upper        fm
	registerStep func(Step)
}

func newRuntaskSubBuilder(upper fm, register func(Step), runner taskRunner) lua.Function {
	rbs := runtaskBuilderState{
		runtaskStep: runtaskStep{
			runner: runner,
		},
		upper:        upper,
		registerStep: register,
	}
	return rbs.runTask
}

func (rbs runtaskBuilderState) runTask(l *lua.State) int {
	otherTaskID := lua.CheckString(l, -1)

	rbs.taskID = otherTaskID

	rbs.registerStep(rbs.runtaskStep)
	return tableWith(l, rbs.upper)
}
