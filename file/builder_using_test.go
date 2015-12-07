package file

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"testing"
)

func ExampleAbsolutizeBinds() {
	h := absolutizeBinds(docker.HostConfig{
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
	defer func() {
		if x := recover(); x == nil {
			panic("Didn't panic")
		}
	}()

	absolutizeBinds(docker.HostConfig{
		Binds: []string{
			"test",
		},
	}, "/projects/alpha")
}
