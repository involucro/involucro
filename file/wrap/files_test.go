package wrap

import (
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestImageConfigFile(t *testing.T) {
	Convey("Given imageid=123 parentid=456", t, func() {
		imageid, parentid := "123", "456"
		Convey("Then the created time is within ten seconds", func() {
			_, buf := imageConfigFile(parentid, imageid, docker.Config{})
			var conf docker.Image
			json.Unmarshal(buf, &conf)
			So(conf.Created, ShouldHappenWithin, time.Duration(10)*time.Second, time.Now())
		})
	})
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
