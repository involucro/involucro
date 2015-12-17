package wrap

import (
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"testing"
	"time"
)

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

func ExampleRepoNameAndTagFrom() {
	repo, tag := repoNameAndTagFrom("foo/bar")
	fmt.Printf("%s %s", repo, tag)
	// Output: foo/bar latest
}

func ExampleRepoNameAndTagFrom_SpecifiedTag() {
	repo, tag := repoNameAndTagFrom("foo/bar:v1")
	fmt.Printf("%s %s", repo, tag)
	// Output: foo/bar v1
}
