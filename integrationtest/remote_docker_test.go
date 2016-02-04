// +build linux

package integrationtest

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/app"
)

func TestRemoteWrappableViaTcp(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	if err := os.MkdirAll("fixture_09", 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("fixture_09/asd", []byte("blahblubb\n"), 0200); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("fixture_09")

	ln, err := net.Listen("tcp4", "127.0.0.1:4243")
	if err != nil {
		t.Fatal(err)
	}

	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				break
			}
			defer c.Close()

			upstream, err := net.Dial("unix", "/var/run/docker.sock")
			if err != nil {
				break
			}

			go func(c, upstream net.Conn) {
				defer upstream.Close()
				io.Copy(c, upstream)
			}(c, upstream)

			go func(c, upstream net.Conn) {
				defer c.Close()
				io.Copy(upstream, c)
			}(c, upstream)
		}
	}()

	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		c.RemoveImage("inttest/9")
	}()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := app.Main([]string{
		"involucro", "-H", "tcp://127.0.0.1:4243", "-l=-2",
		"-w", pwd, "-e",
		"inv.task('wrap').wrap('fixture_09').at('/blah').inImage('busybox').as('inttest/9')",
		"wrap",
	}); err != nil {
		t.Fatal(err)
	}

	_, err = c.InspectImage("inttest/9")
	if err != nil {
		t.Fatal(err)
	}

	assertStdoutContainsFlag([]string{
		"-e", "inv.task('a').using('inttest/9').run('/bin/cat', '/blah/asd')", "a",
	}, "blahblubb", t)
}
