package runtime

import (
	"flag"
	"os"
	"testing"

	"github.com/thriqon/involucro/ilog"
)

func TestMain(m *testing.M) {
	flag.Parse()
	ilog.StdLogger.SetPrintFunc(func(_ ilog.Bough) {})
	os.Exit(m.Run())
}

func newEmpty() *Runtime {
	r := New(make(map[string]string), nil, ".")
	return &r
}
