package integrationtest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/app"
)

func TestWrapNestedDirsCorrectly(t *testing.T) {
	if testing.Short() {
		return
	}
	dir, err := ioutil.TempDir("", "inttest-15")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	nestedPath := filepath.Join(dir, "asd", "p", "aaa")
	if err := os.MkdirAll(nestedPath, 0755); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(nestedPath, "a"), []byte("123"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(nestedPath, "b"), []byte("456"), 0755); err != nil {
		t.Fatal(err)
	}

	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		c.RemoveImage("inttest/15")
	}()

	if err := app.Main([]string{
		"involucro",
		"-e",
		"inv.task('wrap').wrap('.').inImage('busybox').at('/data').as('inttest/15')",
		"-w", dir,
		"wrap",
	}); err != nil {
		t.Fatal(err)
	}

	if err := app.Main([]string{
		"involucro",
		"-e",
		"inv.task('x').using('inttest/15').run('grep', '123', '/data/asd/p/aaa/a').run('grep', '456', '/data/asd/p/aaa/b')",
		"x",
	}); err != nil {
		t.Error(err)
	}
}
