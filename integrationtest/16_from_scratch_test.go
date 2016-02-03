package integrationtest

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/app"
)

func TestFromScratchWithoutParentImage(t *testing.T) {
	if testing.Short() {
		return
	}
	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	args := []string{
		"involucro",
		"-e",
		"inv.task('package').wrap('.').at('/').as('inttest/16')",
		"package",
	}

	defer func() {
		c.RemoveImage("inttest/16")
	}()

	if err := app.Main(args); err != nil {
		t.Fatal(err)
	}

	history, err := c.ImageHistory("inttest/16")
	if err != nil {
		t.Fatal(err)
	}

	if len(history) != 2 {
		t.Fatal("Unexpected history length, was expecting 2", len(history))
	}

	if history[0].Size != 0 {
		t.Error("Top-most layer is not config-only, but has size", history[0].Size)
	}
}
