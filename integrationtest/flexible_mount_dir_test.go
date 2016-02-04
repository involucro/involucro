package integrationtest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/thriqon/involucro/app"
)

func TestFlexibleMountDirs(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pwd, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}

	cases := []string{
		"./int21:/ttt",
		pwd + "/int21:/ttt",
	}

	if err := os.MkdirAll("int21", 0755); err != nil {
		t.Fatal(err)
	}

	for _, el := range cases {
		if err := ioutil.WriteFile(filepath.Join(pwd, "int21", "testfile"), []byte{0}, 0755); err != nil {
			t.Error("case failed", err)
			continue
		}

		if err := app.Main([]string{
			"involucro", "-e",
			"inv.task('p').using('busybox').withHostConfig({Binds = {'" + el + "'}}).run('rm', '/ttt/testfile')",
			"p",
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(filepath.Join(pwd, "int21", "testfile")); (err == nil) || (!os.IsNotExist(err)) {
			t.Error("Unexpected error", err)
		}
	}
}
