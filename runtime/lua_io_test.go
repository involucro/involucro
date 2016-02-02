package runtime

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/Shopify/go-lua"
)

func TestReadLinesFromFile(t *testing.T) {
	t.Parallel()

	f, err := ioutil.TempFile("", "inv-test-io")
	if err != nil {
		t.Fatal("Unable to create temp file", err)
	}
	filename := f.Name()
	defer os.Remove(filename)
	defer f.Close()

	expected := []string{
		"--FIXTURE LINE 1",
		"more data",
		"even windows line ends\r",
		"more data ",
	}

	if _, err := io.WriteString(f, strings.Join(expected, "\n")); err != nil {
		t.Fatal("Unable to write to file", err)
	}
	f.Close()

	l := lua.NewState()
	injectIoLib(l)

	var pos int
	l.Register("push_result", func(l *lua.State) int {
		line := lua.CheckString(l, -1)
		if "LINE:"+strings.TrimSuffix(expected[pos], "\r") != line {
			t.Errorf("Unexpected %s, expected %s", line, expected[pos])
		}
		pos++
		return 0
	})

	if err := lua.DoString(l, `for l in io.lines('`+strings.Replace(filename, "\\", "\\\\", -1)+`') do push_result("LINE:" .. l) end`); err != nil {
		t.Error(err)
	}
}

func TestReadLinesFromNonExistingFile(t *testing.T) {
	t.Parallel()

	l := lua.NewState()
	injectIoLib(l)

	filename := "/non/existing/file.txt"
	l.Register("push_result", func(l *lua.State) int {
		t.Error("result function called")
		panic("unreachable")
	})
	if err := lua.DoString(l, `for l in io.lines('`+strings.Replace(filename, "\\", "\\\\", -1)+`') do push_result("LINE:" .. l) end`); err == nil {
		t.Error("Missing error for non existing file")
		panic("unreachable")
	}
}

func TestReadLinesFromEmptyFile(t *testing.T) {
	t.Parallel()

	l := lua.NewState()
	injectIoLib(l)

	f, err := ioutil.TempFile("", "inv-test-io2")
	if err != nil {
		t.Fatal("Unable to create temp file", err)
	}
	filename := f.Name()
	f.Close()
	defer os.Remove(filename)

	l.Register("push_result", func(l *lua.State) int {
		t.Error("result function called")
		panic("unreachable")
	})
	if err := lua.DoString(l, `for l in io.lines('`+strings.Replace(filename, "\\", "\\\\", -1)+`') do push_result("LINE:" .. l) end`); err != nil {
		t.Error("Error during reading empty file", err)
	}
}
