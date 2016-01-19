package runtime

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func TestRandomFileName(t *testing.T) {
	filename := randomTarballFileName()
	if !strings.Contains(filename, "involucro") {
		t.Errorf("Didn't contain involucro: %s", filename)
	}
	otherFilename := randomTarballFileName()
	if otherFilename == filename {
		t.Errorf("Other filename is not different from the original: %s == %s", otherFilename, filename)
	}

	if _, err := os.Stat(filename); err == nil {
		t.Errorf("Stat succeeded, file shouldn't exist: %s", filename)
	}

	if info, err := os.Stat(filepath.Dir(filename)); err != nil {
		t.Errorf("Parent failed to stat: %s", err)
	} else {
		if !info.IsDir() {
			t.Errorf("Parent should be a directory")
		}
	}
}

func TestImageConfigFileIsGeneratedWithDateWithinTenSeconds(t *testing.T) {
	imageid, parentid := "123", "456"
	_, buf := imageConfigFile(parentid, imageid, docker.Config{})
	var conf docker.Image
	json.Unmarshal(buf, &conf)

	duration := time.Since(conf.Created).Seconds()
	if duration < 0 {
		duration *= -1
	}
	if duration > 10 {
		t.Errorf("Created was more than 10 seconds ago/in more than 10 seconds")
	}
}

func ExampleImageConfigFile_Contents() {
	imageid, parentid := "123", "456"
	_, buf := imageConfigFile(parentid, imageid, docker.Config{})

	var conf docker.Image
	json.Unmarshal(buf, &conf)

	fmt.Println(conf.ID, conf.Parent)
	// Output: 123 456
}

func ExampleImageConfigFile_TarHeader() {
	imageid, parentid := "123", "456"
	header, _ := imageConfigFile(parentid, imageid, docker.Config{})

	fmt.Println(header.Name)
	// Output: 123/json
}

func ExampleRepositoriesFile() {
	_, buf := repositoriesFile("test/gcc:latest", "283028")
	fmt.Printf("%s\n", buf)
	// Output: {"test/gcc":{"latest":"283028"}}
}

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
