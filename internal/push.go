package runtime

import (
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"github.com/involucro/involucro/auth"
	"github.com/involucro/involucro/ilog"
)

type pushStepBuilderState struct {
	pushStep
	upper        fm
	registerStep func(Step)
}

type pushStep struct {
	docker.PushImageOptions
}

func newPushSubBuilder(upper fm, register func(Step)) lua.Function {
	psbs := pushStepBuilderState{
		upper:        upper,
		registerStep: register,
	}
	return psbs.push
}

func (psbs pushStepBuilderState) push(l *lua.State) int {
	opts := &psbs.PushImageOptions
	opts.Name = lua.CheckString(l, 1)
	if l.Top() >= 2 {
		opts.Tag = lua.CheckString(l, 2)
	} else {
		opts.Name, _, opts.Tag = repoNameAndTagFrom(opts.Name)
	}
	psbs.registerStep(psbs.pushStep)

	return tableWith(l, psbs.upper)
}

func (s pushStep) Take(i *Runtime) error {
	ac, foundAuthentication, err := auth.ForServer(serverOfRepo(s.PushImageOptions.Name))
	if err != nil {
		return err
	}

	if err := i.client.PushImage(s.PushImageOptions, ac); err != nil {
		if !foundAuthentication {
			ilog.Warn.Logf("Pull may have failed due to missing authentication information in ~/.involucro")
		}
		return err
	}
	return nil
}

func (s pushStep) ShowStartInfo() {
	logTask.Logf("Push image [%s:%s]", s.Name, s.Tag)
}

func serverOfRepo(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) < 3 {
		return ""
	}
	return strings.Join(parts[:len(parts)-2], "/")
}
