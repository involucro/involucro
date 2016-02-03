package integrationtest

import (
	"os"
	"testing"
)

func TestEnvironmentVariables(t *testing.T) {
	if testing.Short() {
		return
	}

	os.Setenv("INV_MESSAGE", "inv_message")
	defer os.Setenv("INV_MESSAGE", "")

	assertStdoutContainsFlag([]string{
		"-e",
		"inv.task('test').using('busybox').withExpectation({stdout = \"inv_message\"}).run('/bin/echo', ENV.INV_MESSAGE)",
		"test",
	}, "inv_message", t)
}
