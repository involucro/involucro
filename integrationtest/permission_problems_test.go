// +build linux

package integrationtest

import (
	"os"
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/app"
)

func TestCanWrapFilesOnlyReadableForRoot(t *testing.T) {
	if testing.Short() {
		return
	}

	if err := os.MkdirAll("test", 0755); err != nil {
		t.Fatal(err)
	}

	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		app.Main([]string{
			"involucro", "-e", "inv.task('x').using('busybox').run('/bin/sh', '-c', 'rm -f /source/test/only_root')", "x",
		})
		c.RemoveImage("inttest/wrap_root")
	}()

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('x').using('busybox').run('/bin/sh', '-c', 'echo FLAG > /source/test/only_root && chmod 0400 /source/test/only_root')",
		"x",
	}); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open("test/only_root")
	if err == nil {
		file.Close()
		t.Fatal("File opening succeeded")
	}

	if !os.IsPermission(err) {
		t.Fatal("Error was not a permission error, but", err)
	}

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('w').wrap('test').inImage('busybox').at('/data').as('inttest/wrap_root')",
		"w",
	}); err != nil {
		t.Fatal(err)
	}

	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('x').using('inttest/wrap_root').run('cat', '/data/only_root')",
		"x",
	}, "FLAG", t)
}
