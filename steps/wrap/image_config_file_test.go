package wrap

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestImageConfigFile(t *testing.T) {
	Convey("Given imageid=123 parentid=456", t, func() {
		imageid, parentid := "123", "456"
		Convey("Then the buffer contains the expected JSON document", func() {
			_, buf := imageConfigFile(parentid, imageid)
			So(string(buf), ShouldContainSubstring, `"Id":"123"`)
			So(string(buf), ShouldContainSubstring, `"Parent":"456"`)
		})
		Convey("Then the header puts the file into the correct location", func() {
			header, _ := imageConfigFile(parentid, imageid)
			So(header.Name, ShouldResemble, "123/json")
		})
	})
}
