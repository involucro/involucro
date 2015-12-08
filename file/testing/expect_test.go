package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	file "github.com/thriqon/involucro/file"
	run "github.com/thriqon/involucro/file/run"
	"testing"
)

func TestExpectations(t *testing.T) {
	Convey("Given I have an empty state", t, func() {
		inv := file.InstantiateRuntimeEnv(".")

		prefix := `inv.task('asd').using('asd')`

		Convey("Then I can pass in an empty table", func() {
			So(inv.RunString(prefix+`.withExpectation({}).run()`), ShouldBeNil)
		})

		Convey("Then I can specify a task with an expectation", func() {
			inv.RunString(prefix + `.withExpectation({code = 1}).run()`)
		})

		Convey("Then it returns an error when I pass nothing as expectation", func() {
			So(inv.RunString(prefix+`.withExpectation().run()`), ShouldNotBeNil)
		})

		Convey("Then it returns an error when I pass in a number as expectation", func() {
			So(inv.RunString(prefix+`.withExpectation(5).run()`), ShouldNotBeNil)
		})

		Convey("Then it returns an error when I pass two tables as expectation", func() {
			So(inv.RunString(prefix+`.withExpectation({}, {}).run()`), ShouldNotBeNil)
		})

		Convey("Then it returns an error when I pass in a string as expected code", func() {
			So(inv.RunString(prefix+`.withExpectation({code = 'asd'}).run()`), ShouldNotBeNil)
		})

		Convey("Then it returns an error when I pass in a table as expected output", func() {
			So(inv.RunString(prefix+`.withExpectation({stdout = {}}).run()`), ShouldNotBeNil)
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
				Convey("Then it has no ExpectedStdout matcher", func() {
					So(step.ExpectedStdoutMatcher, ShouldBeNil)
				})
				Convey("Then it has no ExpectedStderr matcher", func() {
					So(step.ExpectedStderrMatcher, ShouldBeNil)
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
				Convey("Then it has no ExpectedStdout matcher", func() {
					So(step.ExpectedStdoutMatcher, ShouldBeNil)
				})
				Convey("Then it has no ExpectedStderr matcher", func() {
					So(step.ExpectedStderrMatcher, ShouldBeNil)
				})
			})
		})

		Convey("When I set an expectation to stdout~/asd.../ and stderr~[0-9]*~", func() {
			So(inv.RunString(prefix+`.withExpectation({stdout = "asd...", stderr = "[0-9]*"}).run()`), ShouldBeNil)
			Convey("Then the task has exactly one step", func() {
				So(inv.Tasks["asd"], ShouldHaveLength, 1)
			})
			Convey("And when I retrieve that step", func() {
				step := inv.Tasks["asd"][0].(run.ExecuteImage)
				Convey("Then it validates for exit code 0", func() {
					So(step.ExpectedCode, ShouldEqual, 0)
				})
				Convey("Then it accepts asdasd as stdout", func() {
					So(step.ExpectedStdoutMatcher, shouldAccept, "asdasd")
				})
				Convey("Then it accepts the empty string as stderr", func() {
					So(step.ExpectedStderrMatcher, shouldAccept, "")
				})
				Convey("Then it accepts a number as string as stderr", func() {
					So(step.ExpectedStderrMatcher, shouldAccept, "48304785947")
				})
			})
		})
	})
}
