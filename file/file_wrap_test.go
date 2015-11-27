package file

import (
	. "github.com/smartystreets/goconvey/convey"
	wrap "github.com/thriqon/involucro/steps/wrap"
	"testing"
)

func TestWrapTaskDefinition(t *testing.T) {
	Convey("Given an empty runtime environment", t, func() {
		inv := InstantiateRuntimeEnv(".")
		env := inv.duk

		runCode := func(s string) func() {
			return func() {
				env.EvalString(s)
			}
		}
		Convey("When I define a wrap task", func() {
			runCode(`inv.task('w').wrap("dist").inImage("p").at("/data").as("test/one")`)()
			Convey("Then it has a corresponding entry in the task map", func() {
				So(inv.Tasks["w"], ShouldNotBeEmpty)
				So(inv.Tasks["w"][0], ShouldHaveSameTypeAs, wrap.AsImage{})
				wi := inv.Tasks["w"][0].(wrap.AsImage)
				So(wi.ParentImage, ShouldResemble, "p")
			})
		})
	})
}
