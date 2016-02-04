package integrationtest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/thriqon/involucro/app"
)

func testStdoutOf(f func() error, expected string, t *testing.T) {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdout = stdout
	}()

	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		r.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: copying pipe: %v\n", err)
			os.Exit(1)
		}
		outC <- buf.String()
	}()

	if err := f(); err != nil {
		t.Fatal(err)
	}

	w.Close()
	out := <-outC
	os.Stdout = stdout

	if out != "a\nb\n" {
		t.Errorf("unexpected output %v", out)
	}
}

func TestTaskListWithTAndDirectScript(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	testStdoutOf(func() error {
		return app.Main([]string{
			"involucro", "-e",
			"inv.task('a').using('busybox').run('x'); inv.task('b').using('busybox').run('z')",
			"-T",
		})
	}, "a\nb\n", t)
}
func TestTaskListWithTasksAndDirectScript(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	testStdoutOf(func() error {
		return app.Main([]string{
			"involucro", "-e",
			"inv.task('a').using('busybox').run('x'); inv.task('b').using('busybox').run('z')",
			"--tasks",
		})
	}, "a\nb\n", t)
}
