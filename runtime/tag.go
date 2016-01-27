package runtime

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
)

type tagStep struct {
	originalName    string
	tagImageOptions docker.TagImageOptions
}

type tagStepBuilder struct {
	tagStep
	upper        fm
	registerStep func(Step)
}

func (s tagStep) Take(i *Runtime) error {
	return i.client.TagImage(s.originalName, s.tagImageOptions)
}

// ShowStartInfo displays logging information including the executed task.
func (s tagStep) ShowStartInfo() {
	logTask.Logf("Tag Image [%s] as [%s:%s]", s.originalName, s.tagImageOptions.Repo, s.tagImageOptions.Tag)
}

func newTagSubBuilder(upper fm, register func(Step)) lua.Function {
	tsb := tagStepBuilder{
		upper:        upper,
		registerStep: register,
	}
	return tsb.tag
}

func (tsb tagStepBuilder) tag(l *lua.State) int {
	tsb.originalName = lua.CheckString(l, -1)
	return tableWith(l, tsb.upper, fm{"as": tsb.as})
}

func (tsb tagStepBuilder) as(l *lua.State) int {
	repo, _, tag := repoNameAndTagFrom(lua.CheckString(l, -1))
	// true: force
	tsb.tagImageOptions = docker.TagImageOptions{repo, tag, true}
	tsb.registerStep(tsb.tagStep)
	return tableWith(l, tsb.upper)
}
