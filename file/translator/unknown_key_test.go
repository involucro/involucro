package translator

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"reflect"
	"testing"
)

func TestUnknownPropertiesConfig(t *testing.T) {
	source := `x = {blah = 5}`

	expected := docker.Config{}

	state := lua.NewState()
	if err := lua.DoString(state, source); err != nil {
		t.Errorf("Error executing string: %s", err)
	}
	state.Global("x")

	if actual := ParseImageConfigFromLuaTable(state); !reflect.DeepEqual(actual, expected) {
		t.Errorf("Wasn't unchanged: %s != %s", actual, expected)
	}
}

func TestUnknownPropertiesHostConfig(t *testing.T) {
	source := `x = {blah = 5}`

	expected := docker.HostConfig{}

	state := lua.NewState()
	if err := lua.DoString(state, source); err != nil {
		t.Errorf("Error executing string: %s", err)
	}
	state.Global("x")

	if actual := ParseHostConfigFromLuaTable(state); !reflect.DeepEqual(actual, expected) {
		t.Errorf("Wasn't unchanged: %s != %s", actual, expected)
	}
}
