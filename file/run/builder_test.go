package run

import (
	"fmt"
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"testing"
)

func ExampleAbsolutizeBinds() {
	h, _ := absolutizeBinds(docker.HostConfig{
		Binds: []string{
			"./:/source",
			"/data:/data",
			"dist:/dist",
		},
	}, "/projects/alpha")

	for _, el := range h.Binds {
		fmt.Println(el)
	}
	// Output:
	// /projects/alpha:/source
	// /data:/data
	// /projects/alpha/dist:/dist
}

func TestAbsolutizeBinds(t *testing.T) {
	_, err := absolutizeBinds(docker.HostConfig{
		Binds: []string{
			"test",
		},
	}, "/projects/alpha")
	if err == nil {
		t.Error("Didn't returne error")
	}
}

func ExampleArgumentsToStringArray() {
	l := lua.NewState()
	l.PushString("a")
	l.PushString("s")
	l.PushString("d")
	fmt.Println(argumentsToStringArray(l))
	// Output: [a s d]
}
