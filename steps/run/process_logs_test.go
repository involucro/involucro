package run

import (
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestReadAndMatchAgainst(t *testing.T) {
	Convey("Given I have a blob of data in a reader", t, func() {
		reader := strings.NewReader("asdjsdkfsdjkl")
		Convey("When I use the regex /^asd$/ on it", func() {
			regex := regexp.MustCompile("^asd$")
			Convey("Then readAndMatchAgainst sends back an error", func() {
				ch := make(chan error, 1)
				readAndMatchAgainst(reader, regex, ch, "testing")
				val := <-ch
				So(val, ShouldNotBeNil)
			})
		})
		Convey("When I use the regex /dfd92hj/ on it", func() {
			regex := regexp.MustCompile("dfd92hj")
			Convey("Then readAndMatchAgainst sends back an error", func() {
				ch := make(chan error, 1)
				readAndMatchAgainst(reader, regex, ch, "testing")
				val := <-ch
				So(val, ShouldNotBeNil)
			})
		})
		Convey("When I use the regex /asd.*/ on it", func() {
			regex := regexp.MustCompile("asd.*")
			Convey("Then readAndMatchAgainst accepts that", func() {
				ch := make(chan error, 1)
				readAndMatchAgainst(reader, regex, ch, "testing")
				val := <-ch
				So(val, ShouldBeNil)
			})
		})
	})
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

//func (img *ExecuteImage) loadAndProcessLogs(c dockerLogsProvider, containerID string) error {

func TestProcessLogs(t *testing.T) {
	containerID := "123"
	prov := mockDockerLogsProvider{}

	Convey("Given an ExcuteImage without any matchers", t, func() {
		ei := ExecuteImage{}
		Convey("When asked to process the logs", func() {
			err := ei.loadAndProcessLogs(&prov, containerID)
			Convey("Then it executes without an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then it uses the container ID 123", func() {
				So(prov.lastCalledWith.Container, ShouldResemble, "123")
			})
		})
	})
}
