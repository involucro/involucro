package integrationtest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const literateModeScript = `
# Testing

See, it supports markdown: **emphasis**

We can also define tasks:

> inv.task('blubb')

Altough, interrupted!


I can also use four spaces:

    .using('busybox').run('/bin/echo', 'Test OK')

But it needs to be have an empty row before the code.
    This doesn't crash even though it's illegal Lua.

That's all
`

func TestLiterateModeGivenFilename(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	dir, err := ioutil.TempDir("", "inttest-31")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, "somescript.md")

	if err := ioutil.WriteFile(filename, []byte(literateModeScript), 0755); err != nil {
		t.Fatal(err)
	}

	assertStdoutContainsFlag([]string{
		"-f",
		filename,
		"blubb",
	}, "Test OK", t)
}

func TestLiterateModeAutomaticFilename(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	filename := "invfile.lua.md"
	if err := ioutil.WriteFile(filename, []byte(literateModeScript), 0755); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("invfile.lua.md")

	assertStdoutContainsFlag([]string{
		"blubb",
	}, "Test OK", t)
}
