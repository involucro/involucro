package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	file "github.com/thriqon/involucro/file"
	run "github.com/thriqon/involucro/steps/run"
	"testing"
)

func TestExpectations(t *testing.T) {
	Convey("Given I have an empty state", t, func() {
		inv := file.InstantiateRuntimeEnv(".")

		prefix := `inv.task('asd').using('asd')`

		Convey("Then I can specify a task with an expectation", func() {
			inv.RunString(prefix + `.withExpectation({code = 1}).run()`)
		})

		Convey("Then it panics when I pass nothing as expectation", func() {
			So(inv.RunString(prefix+`.withExpectation().run()`), ShouldNotBeNil)
		})

		Convey("Then it panics when I pass in a number as expectation", func() {
			So(inv.RunString(prefix+`.withExpectation(5).run()`), ShouldNotBeNil)
		})

		Convey("Then it panics when I pass two tables as expectation", func() {
			So(inv.RunString(prefix+`.withExpectation({}, {}).run()`), ShouldNotBeNil)
		})

		Convey("When I don't call withExpectation", func() {
			So(inv.RunString(prefix+`.run()`), ShouldBeNil)
			Convey("Then the task is defined", func() {
				So(inv.Tasks, ShouldContainKey, "asd")
			})
			Convey("Then the task has exactly one step", func() {
				So(inv.Tasks["asd"], ShouldHaveLength, 1)
			})
			Convey("And when I retrieve that step", func() {
				step := inv.Tasks["asd"][0].(run.ExecuteImage)
				Convey("Then it validates for exit code 0 == SUCCESS", func() {
					So(step.ExpectedCode, ShouldEqual, 0)
				})
			})
		})

		Convey("When I set the expected code to 1", func() {
			So(inv.RunString(prefix+`.withExpectation({code = 1}).run()`), ShouldBeNil)
			Convey("Then the task has exactly one step", func() {
				So(inv.Tasks["asd"], ShouldHaveLength, 1)
			})
			Convey("And when I retrieve that step", func() {
				step := inv.Tasks["asd"][0].(run.ExecuteImage)
				Convey("Then it validates for exit code 1", func() {
					So(step.ExpectedCode, ShouldEqual, 1)
				})
			})
		})
	})
}
