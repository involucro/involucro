package runtime

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
)

func TestIsRemoteInstance(t *testing.T) {
	variables := make(map[string]string)

	var client *docker.Client
	var err error

	client, err = docker.NewClient("tcp://blah.de")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	inv := New(variables, client, ".")
	if !inv.isUsingRemoteInstance() {
		t.Error("expected remote instance")
	}

	client, err = docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	inv = New(variables, client, ".")
	if inv.isUsingRemoteInstance() {
		t.Error("expected local instance")
	}
}
