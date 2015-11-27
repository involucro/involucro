package file

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTaskDefinition(t *testing.T) {
	Convey("Given an empty runtime environment", t, func() {
		inv := InstantiateRuntimeEnv(".")
		env := inv.duk

		runCode := func(s string) func() {
			return func() {
				env.EvalString(s)
			}
		}

		Convey("Defining an empty Task succeeds", func() {
			So(runCode(`inv.task('test')`), ShouldNotPanic)
		})

	})
}
