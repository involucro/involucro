
package integrationtest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/involucro/involucro/app"
)

func TestWrapCurrentDir(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	dir, err := ioutil.TempDir("", "inttest-58")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := ioutil.WriteFile(filepath.Join(dir, "a"), []byte("123"), 0755); err != nil {
		t.Fatal(err)
	}

	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		c.RemoveImage("inttest/58")
		os.Chdir(cwd)
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := app.Main([]string{
		"involucro",
		"-e",
		"inv.task('wrap').wrap('.').inImage('busybox').at('/data').as('inttest/15')",
		"wrap",
	}); err != nil {
		t.Fatal(err)
	}

	if err := app.Main([]string{
		"involucro",
		"-e",
		"inv.task('x').using('inttest/15').run('grep', '123', '/data/a')",
		"x",
	}); err != nil {
		t.Error(err)
	}
}
