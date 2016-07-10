package runtime

import (
	"flag"
	"os"
	"testing"

	"github.com/involucro/involucro/ilog"
)

func TestMain(m *testing.M) {
	flag.Parse()
	ilog.StdLog.SetPrintFunc(func(_ ilog.Bough) {})
	os.Exit(m.Run())
}

func newEmpty() *Runtime {
	r := New(make(map[string]string), nil, ".")
	return &r
}
