package file

import (
	"github.com/thriqon/involucro/file/run"
	"testing"
)

func TestExpectations(t *testing.T) {
	t.Parallel()

	inv := InstantiateRuntimeEnv(make(map[string]string))
	prefix := `inv.task('asd').using('asd')`

	if err := inv.RunString(prefix + `.withExpectation({}).run()`); err != nil {
		t.Error(err)
	}

	if err := inv.RunString(prefix + `.withExpectation({code = 1}).run()`); err != nil {
		t.Error(err)
	}

	if err := inv.RunString(prefix + `.withExpectation().run()`); err == nil {
		t.Error("Error expected")
	}

	if inv.RunString(prefix+`.withExpectation(5).run()`) == nil {
		t.Error("Error expected")
	}

	if inv.RunString(prefix+`.withExpectation({}, {}).run()`) == nil {
		t.Error("Error expected")
	}

	if inv.RunString(prefix+`.withExpectation({code = 'asd'}).run()`) == nil {
		t.Error("Error expected for string code")
	}

	if inv.RunString(prefix+`.withExpectation({stdout = {}}).run()`) == nil {
		t.Error("Error expected for table expectation")
	}
}

func TestExpectationsWithoutOverrides(t *testing.T) {
	t.Parallel()

	inv := InstantiateRuntimeEnv(make(map[string]string))
	prefix := `inv.task('asd').using('asd')`
	if err := inv.RunString(prefix + `.run()`); err != nil {
		t.Fatal(err)
	}

	if len(inv.tasks["asd"]) != 1 {
		t.Fatal("Invalid number of steps")
	}

	step := inv.tasks["asd"][0].(run.ExecuteImage)

	if step.ExpectedCode != 0 {
		t.Error("ExpectedCode should be 0, is", step.ExpectedCode)
	}

	if step.ExpectedStdoutMatcher != nil || step.ExpectedStderrMatcher != nil {
		t.Error("Expected output matchers to be nil")
	}

}

func TestExpectationsExpectCode1(t *testing.T) {
	t.Parallel()

	inv := InstantiateRuntimeEnv(make(map[string]string))
	prefix := `inv.task('asd').using('asd')`

	if err := inv.RunString(prefix + `.withExpectation({code = 1}).run()`); err != nil {
		t.Fatal(err)
	}

	if len(inv.tasks["asd"]) != 1 {
		t.Fatal("Invalid number of steps")
	}

	step := inv.tasks["asd"][0].(run.ExecuteImage)

	if step.ExpectedCode != 1 {
		t.Error("ExpectedCode should be 1, is", step.ExpectedCode)
	}

	if step.ExpectedStdoutMatcher != nil || step.ExpectedStderrMatcher != nil {
		t.Error("Expected output matchers to be nil")
	}
}
func TestExpectationsExpectCertainOutputs(t *testing.T) {
	t.Parallel()

	inv := InstantiateRuntimeEnv(make(map[string]string))
	prefix := `inv.task('asd').using('asd')`

	if err := inv.RunString(prefix + `.withExpectation({stdout = "asd...", stderr = "[0-9]*"}).run()`); err != nil {
		t.Fatal(err)
	}

	if len(inv.tasks["asd"]) != 1 {
		t.Fatal("Invalid number of steps")
	}

	step := inv.tasks["asd"][0].(run.ExecuteImage)

	if step.ExpectedCode != 0 {
		t.Error("ExpectedCode should be 0, is", step.ExpectedCode)
	}

	if step.ExpectedStdoutMatcher == nil || step.ExpectedStderrMatcher == nil {
		t.Error("Expected output matchers to not be nil")
	}

	if !step.ExpectedStdoutMatcher.MatchString("asdasd") {
		t.Error("Expected to match asdasd on stdout")
	}
	if !step.ExpectedStderrMatcher.MatchString("") {
		t.Error("Expected to match the empty string on stderr")
	}
	if !step.ExpectedStderrMatcher.MatchString("48304785947") {
		t.Error("Expected to match the given number on stderr")
	}
}
