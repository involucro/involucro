package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	file "github.com/thriqon/involucro/file"
	run "github.com/thriqon/involucro/steps/run"
	"testing"
)

func TestReuseReturnValues(t *testing.T) {
	Convey("Given I have an instance", t, func() {
		inv := file.InstantiateRuntimeEnv(".")
		Convey("When I store the result of .task('test') in a variable", func() {
			So(inv.RunString(`
			test = inv.task('test')
			inv.task('build').using('asd').run('asd')
			test.using('dsa').run('dsa')
			`), ShouldBeNil)

			Convey("Then I have the 'test' task", func() {
				So(inv.Tasks, ShouldContainKey, "test")
				tests := inv.Tasks["test"]
				So(len(tests), ShouldEqual, 1)

				So(tests[0].(run.ExecuteImage).Config.Image, ShouldResemble, "dsa")
				So(tests[0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"dsa"})

			})
			Convey("Then I have the 'builds' task", func() {
				So(inv.Tasks, ShouldContainKey, "build")
				builds := inv.Tasks["build"]
				So(len(builds), ShouldEqual, 1)

				So(builds[0].(run.ExecuteImage).Config.Image, ShouldResemble, "asd")
				So(builds[0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"asd"})
			})
		})
		Convey("When I store the result of .using('dsa') in a variable", func() {
			So(inv.RunString(`
			dsa = inv.task('test').using('dsa')
			inv.task('build').using('asd').run('asd')
			dsa.run('dsa')
			`), ShouldBeNil)
			Convey("Then I have the 'test' task", func() {
				So(inv.Tasks, ShouldContainKey, "test")
				tests := inv.Tasks["test"]
				So(len(tests), ShouldEqual, 1)

				So(tests[0].(run.ExecuteImage).Config.Image, ShouldResemble, "dsa")
				So(tests[0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"dsa"})
			})

			Convey("Then I have the 'builds' task", func() {
				So(inv.Tasks, ShouldContainKey, "build")
				builds := inv.Tasks["build"]
				So(len(builds), ShouldEqual, 1)

				So(builds[0].(run.ExecuteImage).Config.Image, ShouldResemble, "asd")
				So(builds[0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"asd"})
			})
		})
	})
}
