package run

import (
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/file/translator"
	"github.com/thriqon/involucro/file/types"
	"github.com/thriqon/involucro/file/utils"
	"path"
	"regexp"
	"strings"
)

type usingBuilderState struct {
	ExecuteImage
	upper        utils.Fm
	registerStep func(types.Step)
	workingDir   string
}

func NewSubBuilder(upper utils.Fm, register func(types.Step), workingDir string) lua.Function {
	ubs := usingBuilderState{
		workingDir:   workingDir,
		registerStep: register,
		upper:        upper,
		ExecuteImage: ExecuteImage{
			HostConfig: docker.HostConfig{
				Binds: []string{
					"./:/source",
				},
			},
		},
	}

	return ubs.using
}

func (ubs usingBuilderState) usingTable(l *lua.State) int {
	return utils.TableWith(l, ubs.upper, utils.Fm{
		"run":             ubs.usingRun,
		"withExpectation": ubs.usingWithExpectation,
		"withConfig":      ubs.withConfig,
		"withHostConfig":  ubs.withHostConfig,
	})
}

func (ubs usingBuilderState) using(l *lua.State) int {
	ubs.Config.Image = lua.CheckString(l, -1)
	return ubs.usingTable(l)
}

func (ubs usingBuilderState) usingRun(l *lua.State) int {
	ubs.Config.Cmd = argumentsToStringArray(l)
	if ubs.Config.WorkingDir == "" {
		ubs.Config.WorkingDir = "/source"
	}

	ubs.HostConfig = absolutizeBinds(ubs.HostConfig, ubs.workingDir)

	ubs.registerStep(ubs.ExecuteImage)
	return ubs.usingTable(l)
}

func (ubs usingBuilderState) usingWithExpectation(l *lua.State) int {
	if l.Top() != 1 {
		lua.Errorf(l, "expected exactly one argument to 'withExpectation'")
		panic("unreachable")
	}
	lua.ArgumentCheck(l, l.IsTable(-1), 1, "Expected table as argument")

	l.Field(-1, "code")
	if !l.IsNil(-1) {
		ubs.ExpectedCode = lua.CheckInteger(l, -1)
		log.WithFields(log.Fields{"code": ubs.ExpectedCode}).Info("Expecting code")
	}
	l.Pop(1)

	l.Field(-1, "stdout")
	if !l.IsNil(-1) {
		str := lua.CheckString(l, -1)
		if regex, err := regexp.Compile(str); err != nil {
			lua.ArgumentError(l, 1, "invalid regular expression in stdout: "+err.Error())
			panic("unreachable")
		} else {
			ubs.ExpectedStdoutMatcher = regex
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
			ubs.ExpectedStderrMatcher = regex
		}
	}
	l.Pop(1)

	return ubs.usingTable(l)
}

func (ubs usingBuilderState) withConfig(l *lua.State) int {
	oldImageID := ubs.Config.Image
	ubs.Config = translator.ParseImageConfigFromLuaTable(l)
	if ubs.Config.Image != "" {
		log.Warn("Overwriting the used image in withConfig is discouraged")
	} else {
		ubs.Config.Image = oldImageID
	}
	return ubs.usingTable(l)
}

func (ubs usingBuilderState) withHostConfig(l *lua.State) int {
	ubs.HostConfig = translator.ParseHostConfigFromLuaTable(l)
	return ubs.usingTable(l)
}

func absolutizeBinds(h docker.HostConfig, workDir string) docker.HostConfig {
	for ind, el := range h.Binds {
		parts := strings.Split(el, ":")
		if len(parts) != 2 {
			log.WithFields(log.Fields{"bind": el}).Panic("Invalid bind, has to be of the form: source:dest")
		}

		if !path.IsAbs(parts[0]) {
			h.Binds[ind] = path.Join(workDir, parts[0]) + ":" + parts[1]
		}
	}
	return h
}

func argumentsToStringArray(l *lua.State) (args []string) {
	top := l.Top()
	args = make([]string, top)
	for i := 1; i <= top; i++ {
		args[i-1] = lua.CheckString(l, i)
	}
	return
}
