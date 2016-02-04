package integrationtest

import (
	"testing"

	"github.com/thriqon/involucro/app"
)

func TestRuntaskOtherTaskPresent(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('blah').runTask('test'); inv.task('test').using('busybox').run('echo', 'TEST8102')",
		"blah",
	}, "TEST8102", t)
}

func TestRuntaskOtherTaskNotPresent(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	if err := app.Main([]string{
		"involucro",
		"-e",
		"inv.task('test').runTask('udef')",
		"test",
	}); err != nil {
		t.Error(err)
	}
}
