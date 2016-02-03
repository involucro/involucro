package runtime

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
)

func TestDefineTaskWithTableAsIdPanics(t *testing.T) {
	inv := newEmpty()
	if err := inv.RunString(`inv.task({})`); err == nil {
		t.Fatal("Didn't return error")
	}
}

func TestUsingWithoutParameterPanics(t *testing.T) {
	inv := newEmpty()
	if err := inv.RunString(`inv.task('test').using()`); err == nil {
		t.Fatal("Didn't return error")
	}
}

func TestRunWithoutParameterWorks(t *testing.T) {
	inv := newEmpty()
	if err := inv.RunString(`inv.task('test').using('asd').run()`); err != nil {
		t.Fatal(err)
	}
}

func TestUsingRunWithParameterWorks(t *testing.T) {
	inv := newEmpty()
	if err := inv.RunString(`inv.task('test').using('asd').run('test')`); err != nil {
		t.Fatal(err)
	}
}

func TestDefiningRunTask(t *testing.T) {
	inv := newEmpty()
	if err := inv.RunString(`inv.task('test').using('blah').run('test', '123')`); err != nil {
		t.Fatal(err)
	}

	if inv.HasTask("another_task") {
		t.Fatal("Did define the another_task")
	}

	el, ok := inv.tasks["test"]
	if !ok {
		t.Fatal("Didn't define the task test")
	}

	if len(el) == 0 {
		t.Fatal("No step defined in task")
	}

	step := el[0]

	ei := step.(executeImage)

	if !reflect.DeepEqual(ei.Config.Cmd, []string{"test", "123"}) {
		t.Fatal("Didnt store the correct Cmd slice, but: ", ei.Config.Cmd)
	}
}

func testWithParameters(t *testing.T, params ...string) {
	inv := newEmpty()
	paramsQ := make([]string, len(params))
	for i, el := range params {
		paramsQ[i] = "'" + el + "'"
	}
	if err := inv.RunString(`inv.task('test1').using('blah').run(` + strings.Join(paramsQ, ", ") + `)`); err != nil {
		t.Fatal("Failed executing code", params)
	}
	if l := len(inv.tasks["test1"]); l == 0 {
		t.Fatal("No steps for task test1")
	}

	if cmd := inv.tasks["test1"][0].(executeImage).Config.Cmd; !reflect.DeepEqual(cmd, params) {
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
	inv := newEmpty()
	if err := inv.RunString(`inv.task('test').using('blah').run('test').run('test2').using('asd').run('2')`); err != nil {
		t.Fatal("Unable to run code", err)
	}

	if len(inv.tasks["test"]) != 3 {
		t.Fatal("Not exactly three steps in task")
	}
	if image := inv.tasks["test"][0].(executeImage).Config.Image; image != "blah" {
		t.Error("Image is not expected value (blah)", image)
	}
	if cmd := inv.tasks["test"][0].(executeImage).Config.Cmd; !reflect.DeepEqual(cmd, []string{"test"}) {
		t.Error("Cmd is not expected value", cmd)
	}

	if image := inv.tasks["test"][1].(executeImage).Config.Image; image != "blah" {
		t.Error("Image is not expected value (blah)", image)
	}
	if cmd := inv.tasks["test"][1].(executeImage).Config.Cmd; !reflect.DeepEqual(cmd, []string{"test2"}) {
		t.Error("Cmd is not expected value", cmd)
	}

	if image := inv.tasks["test"][2].(executeImage).Config.Image; image != "asd" {
		t.Error("Image is not expected value (asd)", image)
	}
	if cmd := inv.tasks["test"][2].(executeImage).Config.Cmd; !reflect.DeepEqual(cmd, []string{"2"}) {
		t.Error("Cmd is not expected value", cmd)
	}
}

func TestRunTaskDefinitionWithOptions(t *testing.T) {
	inv := newEmpty()

	if err := inv.RunString(`inv.task('test').using('blah').withConfig({ENV = {"FOO=bar"}}).run('test')`); err != nil {
		t.Fatal("Unable to run code", err)
	}
	if len(inv.tasks["test"]) != 1 {
		t.Fatal("test doesn't have exactly one step")
	}
	if env := inv.tasks["test"][0].(executeImage).Config.Env; !reflect.DeepEqual(env, []string{"FOO=bar"}) {
		t.Error("Env has unexpected value", env)
	}
	if image := inv.tasks["test"][0].(executeImage).Config.Image; image != "blah" {
		t.Error("Image has unexpected value", image)
	}

	if err := inv.RunString(`inv.task('test1').using('blah').withConfig({Image = "aaa"}).run('test')`); err != nil {
		t.Fatal("Unable to run code", err)
	}
	if len(inv.tasks["test1"]) != 1 {
		t.Fatal("test1 doesn't have exactly one step")
	}
	if image := inv.tasks["test1"][0].(executeImage).Config.Image; image != "aaa" {
		t.Error("Image not aaa, but", image)
	}
}

func TestAbsolutizeBinds(t *testing.T) {
	h, _ := absolutizeBinds(docker.HostConfig{
		Binds: []string{
			"./:/source",
			"/data:/data",
			"dist:/dist",
		},
	}, "/projects/alpha")

	expected := []string{
		"/projects/alpha:/source",
		"/data:/data",
		"/projects/alpha/dist:/dist",
	}

	for index, el := range h.Binds {
		if expected[index] != el {
			t.Errorf("absolutized bind to %s, but expected to absolutize to %s", el, expected[index])
		}
	}

	_, err := absolutizeBinds(docker.HostConfig{
		Binds: []string{
			"test",
		},
	}, "/projects/alpha")
	if err == nil {
		t.Error("Didn't return error")
	}
}

func TestArgumentsToStringArray(t *testing.T) {
	l := lua.NewState()
	l.PushString("a")
	l.PushString("s")
	l.PushString("d")

	actual := argumentsToStringArray(l)
	if actual[0] != "a" || actual[1] != "s" || actual[2] != "d" {
		t.Errorf("expected [a s d], got %v", actual)
	}
}

type mockDockerLogsProvider struct {
	lastCalledWith docker.LogsOptions
	forStdout      string
	forStderr      string
}

func (md *mockDockerLogsProvider) Logs(l docker.LogsOptions) error {
	io.WriteString(l.OutputStream, md.forStdout)
	io.WriteString(l.ErrorStream, md.forStderr)
	md.lastCalledWith = l
	return nil
}

func TestProcessLogs(t *testing.T) {
	containerID := "123"
	prov := mockDockerLogsProvider{}
	ei := executeImage{}

	if err := ei.loadAndProcessLogs(&prov, containerID); err != nil {
		t.Fatal("Error during load and process", err)
	}
	if x := prov.lastCalledWith.Container; x != "123" {
		t.Error("Unexpected container id", x)
	}
}
