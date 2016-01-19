package runtime

import (
	"flag"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

type NullWriter struct {
}

func (nw NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestMain(m *testing.M) {
	flag.Parse()
	nw := NullWriter{}
	log.SetOutput(nw)
	os.Exit(m.Run())
}

func newEmpty() *Runtime {
	r := New(make(map[string]string), nil, ".")
	return &r
}
