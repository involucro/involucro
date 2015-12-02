package run

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"os"
	"testing"
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
