package integrationtest

import "testing"

func TestRunOptionsCanSetEnvironmentVariables(t *testing.T) {
	if testing.Short() {
		return
	}
	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('a').using('busybox').withConfig({Env = { 'FOO=bar', 'BAZ=baz'}}).run('/bin/sh', '-c', 'echo $FOO $BAZ')",
		"a",
	}, "bar baz", t)
}
