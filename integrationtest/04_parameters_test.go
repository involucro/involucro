package integrationtest

import "testing"

import "github.com/thriqon/involucro/app"

func TestParametersAcceptsInlineScript(t *testing.T) {
	if testing.Short() {
		return
	}

	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('echo').using('busybox').run('echo', 'FLAG_293403')",
		"echo",
	}, "FLAG_293403", t)
}

func TestParametersRejectBothInlineScriptAndFilename(t *testing.T) {
	if testing.Short() {
		return
	}

	if err := app.Main([]string{
		"involucro",
		"-e",
		"inv.task('test').runTask('udef')",
		"-f",
		"custom-invfile.lua",
		"test",
	}); err == nil {
		t.Error("Did not fail")
	}
}
