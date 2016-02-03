package integrationtest

import (
	"testing"

	"github.com/thriqon/involucro/app"
)

func TestExpectationMatchStdout(t *testing.T) {
	if testing.Short() {
		return
	}
	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!'}).run('echo', 'Hello, World!')",
		"test",
	}); err != nil {
		t.Error(err)
	}
}

func TestExpectationsCrashesWhenStdoutNotMet(t *testing.T) {
	if testing.Short() {
		return
	}

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!'}).run('echo', 'Hello, Moon')",
		"test",
	}); err == nil {
		t.Error("Expected error")
	}

}

func TestExpectationMatchStderr(t *testing.T) {
	if testing.Short() {
		return
	}
	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('test').using('busybox').withExpectation({stderr = 'Hello, World'}).run('/bin/sh', '-c', 'echo Hello, World 1>&2')",
		"test",
	}); err != nil {
		t.Error(err)
	}
}

func TestExpectationsCrashesWhenStderrNotMet(t *testing.T) {
	if testing.Short() {
		return
	}

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('test').using('busybox').withExpectation({stderr = 'Hello, World'}).run('/bin/sh', '-c', 'echo Hello, Moon 1>&2')",
		"test",
	}); err == nil {
		t.Error("Expected error")
	}
}

func TestExpectationMatchExitCode(t *testing.T) {
	if testing.Short() {
		return
	}
	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('test').using('busybox').withExpectation({code = 1}).run('false')",
		"test",
	}); err != nil {
		t.Error(err)
	}
}

func TestExpectationsCrashesWhenExitCodeNotMet(t *testing.T) {
	if testing.Short() {
		return
	}

	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('test').using('busybox').withExpectation({code = 1}).run('true')",
		"test",
	}); err == nil {
		t.Error("Expected error")
	}
}
