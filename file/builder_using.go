package file

import (
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	runS "github.com/thriqon/involucro/steps/run"
	"regexp"
)

type usingBuilderState struct {
	builderState
	imageID        string
	expectedCode   int
	expectedStdout *regexp.Regexp
	expectedStderr *regexp.Regexp
}

func (bs builderState) using(l *lua.State) int {
	nbs := usingBuilderState{
		builderState: bs,
		imageID:      requireStringOrFailGracefully(l, -1, "using"),
	}
	return usingTable(l, &nbs)
}

func (ubs usingBuilderState) usingRun(l *lua.State) int {
	args := argumentsToStringArray(l)
	ei := runS.ExecuteImage{
		Config: docker.Config{
			Image:      ubs.imageID,
			Cmd:        args,
			WorkingDir: "/source",
		},
		HostConfig: docker.HostConfig{
			Binds: []string{
				ubs.inv.WorkingDir + ":/source",
			},
		},
		ExpectedCode:          ubs.expectedCode,
		ExpectedStdoutMatcher: ubs.expectedStdout,
		ExpectedStderrMatcher: ubs.expectedStderr,
	}
	tasks := ubs.inv.Tasks
	tasks[ubs.taskID] = append(tasks[ubs.taskID], ei)

	return usingTable(l, &ubs)
}

func usingTable(l *lua.State, ubs *usingBuilderState) int {
	return tableWith(l, fm{
		"using":           ubs.using,
		"run":             ubs.usingRun,
		"task":            ubs.task,
		"withExpectation": ubs.usingWithExpectation,
	})
}

func (ubs usingBuilderState) usingWithExpectation(l *lua.State) int {
	if l.Top() != 1 {
		lua.Errorf(l, "expected exactly one argument to 'withExpectation'")
		panic("unreachable")
	}
	lua.ArgumentCheck(l, l.IsTable(-1), 1, "Expected table as argument")
	nubs := ubs

	l.Field(-1, "code")
	if !l.IsNil(-1) {
		nubs.expectedCode = lua.CheckInteger(l, -1)
		log.WithFields(log.Fields{"code": nubs.expectedCode}).Info("Expecting code")
	}
	l.Pop(1)

	l.Field(-1, "stdout")
	if !l.IsNil(-1) {
		str := lua.CheckString(l, -1)
		if regex, err := regexp.Compile(str); err != nil {
			lua.ArgumentError(l, 1, "invalid regular expression in stdout: "+err.Error())
			panic("unreachable")
		} else {
			nubs.expectedStdout = regex
		}
	}
	l.Pop(1)

	l.Field(-1, "stderr")
	if !l.IsNil(-1) {
		str := lua.CheckString(l, -1)
		if regex, err := regexp.Compile(str); err != nil {
			lua.ArgumentError(l, 1, "invalid regular expression in stderr: "+err.Error())
			panic("unreachable")
		} else {
			nubs.expectedStderr = regex
		}
	}
	l.Pop(1)

	return usingTable(l, &nubs)
}
