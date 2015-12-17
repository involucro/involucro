package testing

import (
	"github.com/thriqon/involucro/file"
	"github.com/thriqon/involucro/file/run"
	"reflect"
	"strings"
	"testing"
)

func TestDefineTaskWithTableAsIdPanics(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))
	if err := inv.RunString(`inv.task({})`); err == nil {
		t.Fatal("Didn't return error")
	}
}

func TestUsingWithoutParameterPanics(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))
	if err := inv.RunString(`inv.task('test').using()`); err == nil {
		t.Fatal("Didn't return error")
	}
}

func TestRunWithoutParameterWorks(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))
	if err := inv.RunString(`inv.task('test').using('asd').run()`); err != nil {
		t.Fatal(err)
	}
}

func TestUsingRunWithParameterWorks(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))
	if err := inv.RunString(`inv.task('test').using('asd').run('test')`); err != nil {
		t.Fatal(err)
	}
}

func TestDefiningRunTask(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))
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

func testWithParameters(t *testing.T, params ...string) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))
	paramsQ := make([]string, len(params))
	for i, el := range params {
		paramsQ[i] = "'" + el + "'"
	}
	if err := inv.RunString(`inv.task('test1').using('blah').run(` + strings.Join(paramsQ, ", ") + `)`); err != nil {
		t.Fatal("Failed executing code", params)
	}
	if l := len(inv.Tasks["test1"]); l == 0 {
		t.Fatal("No steps for task test1")
	}

	if cmd := inv.Tasks["test1"][0].(run.ExecuteImage).Config.Cmd; !reflect.DeepEqual(cmd, params) {
		t.Error("cmd doesn't params: ", cmd, params)
	}
}

func TestRunTaskDefinitionMultipleParameterLengths(t *testing.T) {
	testWithParameters(t, "test")
	testWithParameters(t, "test", "asd")
	testWithParameters(t, "test", "asd", "123")
	testWithParameters(t, "test", "asd", "123", "dsa")
}

func TestRunTaskDefinitionMultipleSteps(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))
	if err := inv.RunString(`inv.task('test').using('blah').run('test').run('test2').using('asd').run('2')`); err != nil {
		t.Fatal("Unable to run code", err)
	}

	if len(inv.Tasks["test"]) != 3 {
		t.Fatal("Not exactly three steps in task")
	}
	if image := inv.Tasks["test"][0].(run.ExecuteImage).Config.Image; image != "blah" {
		t.Error("Image is not expected value (blah)", image)
	}
	if cmd := inv.Tasks["test"][0].(run.ExecuteImage).Config.Cmd; !reflect.DeepEqual(cmd, []string{"test"}) {
		t.Error("Cmd is not expected value", cmd)
	}

	if image := inv.Tasks["test"][1].(run.ExecuteImage).Config.Image; image != "blah" {
		t.Error("Image is not expected value (blah)", image)
	}
	if cmd := inv.Tasks["test"][1].(run.ExecuteImage).Config.Cmd; !reflect.DeepEqual(cmd, []string{"test2"}) {
		t.Error("Cmd is not expected value", cmd)
	}

	if image := inv.Tasks["test"][2].(run.ExecuteImage).Config.Image; image != "asd" {
		t.Error("Image is not expected value (asd)", image)
	}
	if cmd := inv.Tasks["test"][2].(run.ExecuteImage).Config.Cmd; !reflect.DeepEqual(cmd, []string{"2"}) {
		t.Error("Cmd is not expected value", cmd)
	}
}

func TestRunTaskDefinitionWithOptions(t *testing.T) {
	inv := file.InstantiateRuntimeEnv(make(map[string]string))

	if err := inv.RunString(`inv.task('test').using('blah').withConfig({ENV = {"FOO=bar"}}).run('test')`); err != nil {
		t.Fatal("Unable to run code", err)
	}
	if len(inv.Tasks["test"]) != 1 {
		t.Fatal("test doesn't have exactly one step")
	}
	if env := inv.Tasks["test"][0].(run.ExecuteImage).Config.Env; !reflect.DeepEqual(env, []string{"FOO=bar"}) {
		t.Error("Env has unexpected value", env)
	}
	if image := inv.Tasks["test"][0].(run.ExecuteImage).Config.Image; image != "blah" {
		t.Error("Image has unexpected value", image)
	}

	if err := inv.RunString(`inv.task('test1').using('blah').withConfig({Image = "aaa"}).run('test')`); err != nil {
		t.Fatal("Unable to run code", err)
	}
	if len(inv.Tasks["test1"]) != 1 {
		t.Fatal("test1 doesn't have exactly one step")
	}
	if image := inv.Tasks["test1"][0].(run.ExecuteImage).Config.Image; image != "aaa" {
		t.Error("Image not aaa, but", image)
	}
}
