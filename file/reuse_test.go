package file

import (
	"github.com/thriqon/involucro/file/run"
	"testing"
)

func TestReuseReturnValuesStoreTaskTest(t *testing.T) {
	t.Parallel()

	inv := InstantiateRuntimeEnv(make(map[string]string))
	if err := inv.RunString(`test = inv.task('test'); inv.task('build').using('asd').run('asd'); test.using('dsa').run('dsa')`); err != nil {
		t.Fatal("Err not nil", err)
	}

	if !inv.HasTask("test") {
		t.Fatal("test task does not exit")
	}

	test := inv.tasks["test"]
	if len(test) != 1 {
		t.Fatal("Invalid number of steps")
	}

	config := test[0].(run.ExecuteImage).Config
	if config.Image != "dsa" || len(config.Cmd) != 1 || config.Cmd[0] != "dsa" {
		t.Fatal("Unexpected configuration values", config)
	}

	if !inv.HasTask("build") {
		t.Fatal("build task does not exist")
	}

	build := inv.tasks["build"]
	if len(build) != 1 {
		t.Fatal("Invalid number of steps")
	}

	config = build[0].(run.ExecuteImage).Config
	if config.Image != "asd" || len(config.Cmd) != 1 || config.Cmd[0] != "asd" {
		t.Error("Unexpected configuration values")
	}

}

func TestReuseReturnValuesStoreTaskDsa(t *testing.T) {
	t.Parallel()

	inv := InstantiateRuntimeEnv(make(map[string]string))
	if err := inv.RunString(`dsa = inv.task('test').using('dsa'); inv.task('build').using('asd').run('asd'); dsa.run('dsa')`); err != nil {
		t.Fatal("err not nil", nil)
	}

	if !inv.HasTask("test") {
		t.Fatal("test task does not exit")
	}

	test := inv.tasks["test"]
	if len(test) != 1 {
		t.Fatal("Invalid number of steps")
	}

	config := test[0].(run.ExecuteImage).Config
	if config.Image != "dsa" || len(config.Cmd) != 1 || config.Cmd[0] != "dsa" {
		t.Fatal("Unexpected configuration values", config)
	}

	if !inv.HasTask("build") {
		t.Fatal("build task does not exist")
	}

	build := inv.tasks["build"]
	if len(build) != 1 {
		t.Fatal("Invalid number of steps")
	}

	config = build[0].(run.ExecuteImage).Config
	if config.Image != "asd" || len(config.Cmd) != 1 || config.Cmd[0] != "asd" {
		t.Error("Unexpected configuration values")
	}
}
