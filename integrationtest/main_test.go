package integrationtest

import (
	"os"
	"testing"

	"github.com/involucro/involucro/ilog"
)

func TestMain(m *testing.M) {
	if !testing.Verbose() {
		ilog.StdLog.SetPrintFunc(func(b ilog.Bough) {})
	}
	os.Exit(m.Run())
}
