package wrap

import (
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestImageConfigFile(t *testing.T) {
	Convey("Given imageid=123 parentid=456", t, func() {
		imageid, parentid := "123", "456"
		Convey("Then the buffer contains the expected JSON document", func() {
			_, buf := imageConfigFile(parentid, imageid, docker.Config{})
			So(string(buf), ShouldContainSubstring, `"Id":"123"`)
			So(string(buf), ShouldContainSubstring, `"Parent":"456"`)
		})
		Convey("Then the created time is within ten seconds", func() {
			_, buf := imageConfigFile(parentid, imageid, docker.Config{})
			var conf docker.Image
			json.Unmarshal(buf, &conf)
			So(conf.Created, ShouldHappenWithin, time.Duration(10)*time.Second, time.Now())
		})
		Convey("Then the header puts the file into the correct location", func() {
			header, _ := imageConfigFile(parentid, imageid, docker.Config{})
			So(header.Name, ShouldResemble, "123/json")
		})
	})
}
