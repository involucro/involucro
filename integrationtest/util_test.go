package integrationtest

import (
	"runtime/debug"
	"testing"

	"github.com/thriqon/involucro/app"
	"github.com/thriqon/involucro/ilog"
)

func assertStdoutContainsFlag(args []string, lineFlag string, t *testing.T) {
	oldPrint := ilog.StdLogger.PrintFunc()
	defer ilog.StdLogger.SetPrintFunc(oldPrint)

	args = append([]string{"involucro", "-l=-2"}, args...)

	var found bool
	ilog.StdLogger.SetPrintFunc(func(b ilog.Bough) {
		if testing.Verbose() && oldPrint != nil {
			oldPrint(b)
		}
		if b.Prefix == "SOUT" && b.Message == lineFlag {
			found = true
		}
	})

	if err := app.Main(args); err != nil {
		debug.PrintStack()
		t.Fatal(err)
	}

	if !found {
		t.Error("Did not find expected flag", lineFlag)
	}
}
