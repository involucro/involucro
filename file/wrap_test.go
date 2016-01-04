package file

import (
	"github.com/thriqon/involucro/file/wrap"
	"testing"
)

func TestWrapTaskDefinition(t *testing.T) {
	inv := InstantiateRuntimeEnv(make(map[string]string))

	if err := inv.RunString(`inv.task('w').wrap("dist").inImage("p").at("/data").as("test/one")`); err != nil {
		t.Fatal("Unable to run code", err)
	}
	if _, ok := inv.tasks["w"]; !ok {
		t.Fatal("w not present as task")
	}
	if len(inv.tasks["w"]) == 0 {
		t.Fatal("w has no steps")
	}
	if _, ok := inv.tasks["w"][0].(wrap.AsImage); !ok {
		t.Fatal("Step is of wrong type")
	}
	if wi := inv.tasks["w"][0].(wrap.AsImage); wi.ParentImage != "p" {
		t.Error("Parent image is unexpected", wi.ParentImage)
	}
}
