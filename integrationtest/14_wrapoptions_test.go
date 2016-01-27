package integrationtest

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/app"
)

func TestWrapOptionSetEntrypoint(t *testing.T) {
	if testing.Short() {
		return
	}

	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		c.RemoveImage("inttest/14")
	}()

	cases := []string{
		"inv.task('wrap').wrap('.').inImage('busybox').at('/data').withConfig({Entrypoint = {'/bin/echo', 'Hello_Options'}}).as('inttest/14')",
		"inv.task('wrap').wrap('.').at('/data').withConfig({Entrypoint = {'/bin/echo', 'Hello_Options'}}).as('inttest/14')",
		"inv.task('wrap').wrap('.').inImage('alpine').at('/data').withConfig({Entrypoint = {'/bin/echo', 'Hello_Options'}}).as('inttest/14')",
	}

	for index, el := range cases {
		if err := app.Main([]string{"involucro", "-e", el, "wrap"}); err != nil {
			t.Error(err)
		}

		image, err := c.InspectImage("inttest/14")
		if err != nil {
			t.Fatalf("Test case %v failed with %v", index, err)
		}
		if image == nil {
			t.Error("Image is nil")
		}
		if image.Config == nil {
			t.Error("Image Config is nil")
		}
		if len(image.Config.Entrypoint) != 2 {
			t.Error("Unexpected entrypoint", image.Config.Entrypoint)
		}
		if image.Config.Entrypoint[0] != "/bin/echo" {
			t.Error("Unexpected first entrypoint element", image.Config.Entrypoint[0])
		}
		if image.Config.Entrypoint[1] != "Hello_Options" {
			t.Error("Unexpected first entrypoint element", image.Config.Entrypoint[0])
		}

		if err := c.RemoveImage("inttest/14"); err != nil {
			t.Error("Unable to remove image", err)
		}
	}
}
