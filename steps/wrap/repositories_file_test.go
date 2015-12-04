package wrap

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRepositoriesFile(t *testing.T) {
	Convey("Given repository=test/helloworld, tag=v1 and image ID 3739474", t, func() {
		repository, imageid := "test/helloworld:v1", "3739474"
		Convey("When calculating the repositories file", func() {
			_, buf := repositoriesFile(repository, imageid)
			Convey("Then the buffer contains the expected JSON document", func() {
				So(string(buf), ShouldResemble, `{"test/helloworld":{"v1":"3739474"}}`)
			})
		})
	})
}

func ExampleRepositoriesFile() {
	_, buf := repositoriesFile("test/gcc:latest", "283028")
	fmt.Printf("%s\n", buf)
	// Output: {"test/gcc":{"latest":"283028"}}
}

func TestRepositoryNameParsing(t *testing.T) {
	Convey("test/asd => [test/asd latest]", t, func() {
		repo, tag := repoNameAndTagFrom("test/asd")
		So(repo, ShouldResemble, "test/asd")
		So(tag, ShouldResemble, "latest")
	})
	Convey("test/asd:v1 => [test/asd v1]", t, func() {
		repo, tag := repoNameAndTagFrom("test/asd:v1")
		So(repo, ShouldResemble, "test/asd")
		So(tag, ShouldResemble, "v1")
	})
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
