package integrationtest

import "testing"

func TestVariablesWithS(t *testing.T) {
	if testing.Short() {
		return
	}
	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('test').using('busybox').run('/bin/echo', VAR['k'])",
		"-s", "k=asd",
		"test",
	}, "asd", t)
}

func TestVariablesWithSet(t *testing.T) {
	if testing.Short() {
		return
	}
	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('test').using('busybox').run('/bin/echo', VAR['k'])",
		"--set", "k=asd",
		"test",
	}, "asd", t)
}

func TestVariablesWithSetAndAdditionalEqualSign(t *testing.T) {
	if testing.Short() {
		return
	}
	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('test').using('busybox').run('/bin/echo', VAR['k'])",
		"--set", "k=asd=6",
		"test",
	}, "asd=6", t)
}
