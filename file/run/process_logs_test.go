package run

import (
	"github.com/fsouza/go-dockerclient"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestReadAndMatchAgainst(t *testing.T) {
	ch := make(chan error, 1)
	readAndMatchAgainst(strings.NewReader("asdjsdkfsdjkl"), regexp.MustCompile("^asd$"), ch, "testing")
	if val := <-ch; val == nil {
		t.Error("Expected error")
	}

	ch = make(chan error, 1)
	readAndMatchAgainst(strings.NewReader("asdjsdkfsdjkl"), regexp.MustCompile("dfd92hj"), ch, "testing")
	if val := <-ch; val == nil {
		t.Error("Unexpectedly no error")
	}

	ch = make(chan error, 1)
	readAndMatchAgainst(strings.NewReader("asdjsdkfsdjkl"), regexp.MustCompile("asd.*"), ch, "testing")
	if val := <-ch; val != nil {
		t.Error("Unexpected error", val)
	}
}

type mockDockerLogsProvider struct {
	lastCalledWith docker.LogsOptions
	forStdout      string
	forStderr      string
}

func (md *mockDockerLogsProvider) Logs(l docker.LogsOptions) error {
	io.WriteString(l.OutputStream, md.forStdout)
	io.WriteString(l.ErrorStream, md.forStderr)
	md.lastCalledWith = l
	return nil
}

func TestProcessLogs(t *testing.T) {
	containerID := "123"
	prov := mockDockerLogsProvider{}
	ei := ExecuteImage{}

	if err := ei.loadAndProcessLogs(&prov, containerID); err != nil {
		t.Fatal("Error during load and process", err)
	}
	if x := prov.lastCalledWith.Container; x != "123" {
		t.Error("Unexpected container id", x)
	}
}
