package runtime

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestPackItUpPrepared(t *testing.T) {
	dir, err := ioutil.TempDir("", "involucro-test-wrap-packitup")
	if err != nil {
		t.Fatal("Unable to create temp dir", err)
	}
	defer os.RemoveAll(dir)

	os.MkdirAll(path.Join(dir, "a", "b", "c", "d"), 0777)
	ioutil.WriteFile(path.Join(dir, "a", "b", "asd"), []byte("ASD"), 0777)

	var buf bytes.Buffer
	if err := packItUp(dir, &buf, "blubb"); err != nil {
		t.Fatal("Unable to pack it up", err)
	}

	tarReader := tar.NewReader(&buf)
	var contents []byte
	expected := map[string]int{
		"blubb":         1,
		"blubb/a":       1,
		"blubb/a/b":     1,
		"blubb/a/b/asd": 1,
		"blubb/a/b/c":   1,
		"blubb/a/b/c/d": 1,
	}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		delete(expected, header.Name)
		if header.Name == "blubb/a/b/asd" {
			var err error
			contents, err = ioutil.ReadAll(tarReader)
			if err != nil {
				t.Fatal("Error during reading", err)
			}
		}
	}

	if fmt.Sprint(expected) != fmt.Sprint(map[string]int{}) {
		t.Error("Didn't see all expected keys, missing were", expected)
	}

	if string(contents) != "ASD" {
		t.Errorf("Contents didn't receive the expected value, were %s", string(contents))
	}
}

func TestPackItUpNotPreparedDir(t *testing.T) {
	var buf bytes.Buffer
	err := packItUp("/not_existing", &buf, "blubb")
	if err == nil {
		t.Error("Didn't return an error when accessing not existent directory")
	}
}

func TestTarHeaderPrepare(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		oldbase, newbase, path, expected string
	}{
		{"p", "x", "p/a/b/c", "x/a/b/c"},
		{"y", "x", "p/a/b/c", "p/a/b/c"},
		{wd, "o", filepath.Join(wd, "dir", "2.txt"), "o/dir/2.txt"},
	}

	for _, el := range cases {
		if actual := preparePathForTarHeader(el.path, el.oldbase, el.newbase); actual != el.expected {
			t.Errorf("expected %s to be rebased to %s, but was to %s", el.path, el.expected, actual)
		}
	}
}

func TestPreparePathForTarHeader(t *testing.T) {
	expected := "asd/compiled/dist/a"
	actual := preparePathForTarHeader("/target/compiled/dist/a", "/target/", "/asd")
	if expected != actual {
		t.Errorf("[%s] is not equal to expected [%s]", actual, expected)
	}
}

func TestWrapTaskDefinition(t *testing.T) {
	inv := newEmpty()

	if err := inv.RunString(`inv.task('w').wrap("dist").inImage("p").at("/data").as("test/one")`); err != nil {
		t.Fatal("Unable to run code", err)
	}
	if _, ok := inv.tasks["w"]; !ok {
		t.Fatal("w not present as task")
	}
	if len(inv.tasks["w"]) == 0 {
		t.Fatal("w has no steps")
	}
	if _, ok := inv.tasks["w"][0].(asImage); !ok {
		t.Fatal("Step is of wrong type")
	}
	if wi := inv.tasks["w"][0].(asImage); wi.ParentImage != "p" {
		t.Error("Parent image is unexpected", wi.ParentImage)
	}
}
