package integrationtest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/involucro/involucro/app"
)

func TestFileHandlingSymlinks(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	dir, err := ioutil.TempDir("", "inttest-26")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := os.Mkdir(filepath.Join(dir, "p"), 0755); err != nil {
		t.Fatal(err)
	}

	if err, ok := os.Symlink("../p", filepath.Join(dir, "p", "cur")).(*os.LinkError); ok && err != nil {
		t.Fatal(err)
	}

	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	args := []string{
		"involucro",
		"-e",
		"inv.task('wrap').wrap('" + dir + "').at('/data').inImage('busybox').as('inttest/26')",
		"wrap",
	}
	defer func() {
		c.RemoveImage("inttest/26")
	}()

	if err := app.Main(args); err != nil {
		t.Fatal(err)
	}

	args = []string{
		"involucro",
		"-e",
		"inv.task('test').using('inttest/26').run('ls', '/data/p/cur/cur/cur/cur')",
		"test",
	}

	if err := app.Main(args); err != nil {
		t.Error(err)
	}
}
