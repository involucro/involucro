package integrationtest

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/app"
)

func TestTagging(t *testing.T) {
	if testing.Short() {
		return
	}
	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		c.RemoveImage("inttest/20:v1")
		c.RemoveImage("inttest/20_b")
	}()

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('package').wrap('.').inImage('busybox').at('/data').as('inttest/20:v1').tag('inttest/20:v1').as('inttest/20_b')",
		"package",
	}); err != nil {
		t.Fatal(err)
	}

	image, err := c.InspectImage("inttest/20:v1")
	if err != nil {
		t.Fatal(err)
	}

	image2, err := c.InspectImage("inttest/20_b")
	if err != nil {
		t.Fatal(err)
	}

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('run').using('inttest/20:v1').run('true')",
		"run",
	}); err != nil {
		t.Fatal(err)
	}

	if image.ID != image2.ID {
		t.Error("Images do not share an ID", image.ID, image2.ID)
	}
}
