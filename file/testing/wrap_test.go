package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	file "github.com/thriqon/involucro/file"
	wrap "github.com/thriqon/involucro/file/wrap"
	"testing"
)

func TestWrapTaskDefinition(t *testing.T) {
	Convey("Given an empty runtime environment", t, func() {
		inv := file.InstantiateRuntimeEnv(".")

		Convey("When I define a wrap task", func() {
			err := inv.RunString(`inv.task('w').wrap("dist").inImage("p").at("/data").as("test/one")`)
			So(err, ShouldBeNil)
			Convey("Then it has a corresponding entry in the task map", func() {
				So(inv.Tasks, ShouldContainKey, "w")
				So(inv.Tasks["w"], ShouldNotBeEmpty)
				So(inv.Tasks["w"][0], ShouldHaveSameTypeAs, wrap.AsImage{})

				wi := inv.Tasks["w"][0].(wrap.AsImage)
				So(wi.ParentImage, ShouldResemble, "p")
			})
		})
	})
}
