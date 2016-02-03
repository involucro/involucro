package integrationtest

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/app"
)

func TestAutopullImage(t *testing.T) {
	if testing.Short() {
		return
	}
	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	c.RemoveImage("tianon/true")

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('test').using('tianon/true').run()",
		"test",
	}); err != nil {
		t.Error(err)
	}
}
