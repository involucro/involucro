package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	file "github.com/thriqon/involucro/file"
	run "github.com/thriqon/involucro/steps/run"
	"testing"
)

func TestRunTaskDefinition(t *testing.T) {
	Convey("Given an empty runtime environment", t, func() {
		inv := file.InstantiateRuntimeEnv(".")

		runCode := func(s string) func() {
			return func() {
				err := inv.RunString(s)
				if err != nil {
					panic(err)
				}
			}
		}

		Convey("Defining a Task with using/run succeeds", func() {
			So(runCode(`inv.task('test').using('blah').run('test')`), ShouldNotPanic)
		})

		Convey("Defining a Task with a table as ID panics", func() {
			So(runCode(`inv.task({})`), ShouldPanic)
		})

		Convey("Calling using without a parameter panics", func() {
			So(runCode(`inv.task('test').using()`), ShouldPanic)
		})

		Convey("Calling run without a parameter does not panic", func() {
			So(runCode(`inv.task('test').using('blah').run()`), ShouldNotPanic)
		})

		Convey("When defining a task with one using/run", func() {
			So(inv.RunString(`inv.task('test').using('blah').run('test', '123')`), ShouldBeNil)

			Convey("Then the task map has not an entry for another task", func() {
				_, ok := inv.Tasks["another_task"]
				So(ok, ShouldBeFalse)
			})

			Convey("Then the task map has an entry for that task", func() {
				_, ok := inv.Tasks["test"]
				So(ok, ShouldBeTrue)
			})

			Convey("Then the task map entry is an ExecuteImage struct", func() {
				So(len(inv.Tasks["test"]), ShouldBeGreaterThan, 0)
				So(inv.Tasks["test"][0], ShouldHaveSameTypeAs, run.ExecuteImage{})
			})

			Convey("Then the task map entry has the given command set", func() {
				So(len(inv.Tasks["test"]), ShouldBeGreaterThan, 0)
				ei := inv.Tasks["test"][0].(run.ExecuteImage)
				So(ei.Config.Cmd, ShouldResemble, []string{"test", "123"})
			})
		})
		Convey("Passing arguments to run works with different lengths", func() {
			So(inv.RunString(`inv.task('test1').using('blah').run('test')`), ShouldBeNil)
			So(inv.RunString(`inv.task('test2').using('blah').run('test', 'asd')`), ShouldBeNil)
			So(inv.RunString(`inv.task('test3').using('blah').run('test', 'asd', '123')`), ShouldBeNil)
			So(inv.RunString(`inv.task('test4').using('blah').run('test', 'asd', '123', 'dsa')`), ShouldBeNil)

			So(len(inv.Tasks["test1"]), ShouldBeGreaterThan, 0)
			So(len(inv.Tasks["test2"]), ShouldBeGreaterThan, 0)
			So(len(inv.Tasks["test3"]), ShouldBeGreaterThan, 0)
			So(len(inv.Tasks["test4"]), ShouldBeGreaterThan, 0)
			So(inv.Tasks["test1"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test"})
			So(inv.Tasks["test2"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test", "asd"})
			So(inv.Tasks["test3"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test", "asd", "123"})
			So(inv.Tasks["test4"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test", "asd", "123", "dsa"})
		})

		Convey("Executing multiple run Tasks after each other works", func() {
			So(inv.RunString(`inv.task('test').using('blah').run('test').run('test2').using('asd').run('2')`), ShouldBeNil)
			So(len(inv.Tasks["test"]), ShouldBeGreaterThan, 2)
			So(inv.Tasks["test"][0].(run.ExecuteImage).Config.Image, ShouldEqual, "blah")
			So(inv.Tasks["test"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test"})

			So(inv.Tasks["test"][1].(run.ExecuteImage).Config.Image, ShouldEqual, "blah")
			So(inv.Tasks["test"][1].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test2"})

			So(inv.Tasks["test"][2].(run.ExecuteImage).Config.Image, ShouldEqual, "asd")
			So(inv.Tasks["test"][2].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"2"})
		})
	})
}
