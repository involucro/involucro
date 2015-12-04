package file

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

func ExampleAbsolutizeBinds() {
	h := docker.HostConfig{
		Binds: []string{
			"./:/source",
			"/data:/data",
			"dist:/dist",
		},
	}
	h2 := absolutizeBinds(h, "/projects/alpha")

	for _, el := range h2.Binds {
		fmt.Println(el)
	}
	// Output:
	// /projects/alpha:/source
	// /data:/data
	// /projects/alpha/dist:/dist
}
