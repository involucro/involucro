package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/thriqon/involucro/file"
	"github.com/thriqon/involucro/file/run"
	"reflect"
	"testing"
)

func TestDefineTaskWithTableAsIdPanics(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(".")
	if err := inv.RunString(`inv.task({})`); err == nil {
		t.Fatal("Didn't return error")
	}
}

func TestUsingWithoutParameterPanics(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(".")
	if err := inv.RunString(`inv.task('test').using()`); err == nil {
		t.Fatal("Didn't return error")
	}
}

func TestRunWithoutParameterWorks(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(".")
	if err := inv.RunString(`inv.task('test').using('asd').run()`); err != nil {
		t.Fatal(err)
	}
}

func TestUsingRunWithParameterWorks(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(".")
	if err := inv.RunString(`inv.task('test').using('asd').run('test')`); err != nil {
		t.Fatal(err)
	}
}

func TestDefiningRunTask(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(".")
	if err := inv.RunString(`inv.task('test').using('blah').run('test', '123')`); err != nil {
		t.Fatal(err)
	}

	if _, ok := inv.Tasks["another_task"]; ok {
		t.Fatal("Did define the another_task")
	}

	el, ok := inv.Tasks["test"]
	if !ok {
		t.Fatal("Didn't define the task test")
	}

	if len(el) == 0 {
		t.Fatal("No step defined in task")
	}

	step := el[0]

	ei := step.(run.ExecuteImage)

	if !reflect.DeepEqual(ei.Config.Cmd, []string{"test", "123"}) {
		t.Fatal("Didnt store the correct Cmd slice, but: ", ei.Config.Cmd)
	}
}

func TestRunTaskDefinition(t *testing.T) {
	Convey("Given an empty runtime environment", t, func() {
		inv := file.InstantiateRuntimeEnv(".")

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

		Convey("When I ask for options", func() {
			So(inv.RunString(`inv.task('test').using('blah').withConfig({ENV = {"FOO=bar"}}).run('test')`), ShouldBeNil)
			So(inv.Tasks["test"], ShouldHaveLength, 1)
			Convey("Then it has that option set in the resulting container config", func() {
				So(inv.Tasks["test"][0].(run.ExecuteImage).Config.Env, ShouldResemble, []string{"FOO=bar"})
			})
			Convey("Then it has blah as image id", func() {
				So(inv.Tasks["test"][0].(run.ExecuteImage).Config.Image, ShouldResemble, "blah")
			})
		})
		Convey("When I overwrite the image id in withConfig", func() {
			So(inv.RunString(`inv.task('test').using('blah').withConfig({Image = "aaa"}).run('test')`), ShouldBeNil)
			So(inv.Tasks["test"], ShouldHaveLength, 1)
			Convey("Then it has aaa as image id", func() {
				So(inv.Tasks["test"][0].(run.ExecuteImage).Config.Image, ShouldResemble, "aaa")
			})
		})
	})
}
