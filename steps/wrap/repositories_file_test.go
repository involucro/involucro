package wrap

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRepositoriesFile(t *testing.T) {
	Convey("Given repository=test/helloworld, tag=v1 and image ID 3739474", t, func() {
		repository, tag, imageid := "test/helloworld", "v1", "3739474"
		Convey("When calculating the repositories file", func() {
			_, buf := repositoriesFile(repository, tag, imageid)
			Convey("Then the buffer contains the expected JSON document", func() {
				So(string(buf), ShouldResemble, `{"test/helloworld":{"v1":"3739474"}}`)
			})
		})
	})
}

func ExampleRepositoriesFile() {
	_, buf := repositoriesFile("test/gcc", "latest", "283028")
	fmt.Printf("%s\n", buf)
	// Output: {"test/gcc":{"latest":"283028"}}
}
