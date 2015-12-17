package wrap

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPackItUpPrepared(t *testing.T) {
	dir, err := ioutil.TempDir("", "involucro-test-wrap-packitup")
	if err != nil {
		t.Fatal("Unable to create temp dir", err)
	}
	defer os.RemoveAll(dir)

	os.MkdirAll(filepath.Join(dir, "a", "b", "c", "d"), 0777)
	ioutil.WriteFile(filepath.Join(dir, "a", "b", "asd"), []byte("ASD"), 0777)

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

func ExampleRebaseFilename() {
	fmt.Println(rebaseFilename("p", "x", "p/a/b/c"))
	fmt.Println(rebaseFilename("y", "x", "p/a/b/c"))
	// Output:
	// x/a/b/c
	// p/a/b/c
}

func ExamplePreparePathForTarHeader() {
	fmt.Println(preparePathForTarHeader("/target/compiled/dist/a", "/target/", "/asd"))
	// Output: asd/compiled/dist/a
}
