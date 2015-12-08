package translator

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnknownPropertiesConfig(t *testing.T) {
	Convey("When I try to set an unknown key", t, func() {
		source := `x = {blah = 5}`

		expected := docker.Config{}

		state := lua.NewState()
		So(lua.DoString(state, source), ShouldBeNil)
		state.Global("x")

		Convey("Then it is accepted and results in an unchaged result", func() {
			actual := ParseImageConfigFromLuaTable(state)
			So(actual, ShouldResemble, expected)
		})
	})
}

func TestUnknownPropertiesHostConfig(t *testing.T) {
	Convey("When I try to set an unknown key", t, func() {
		source := `x = {blah = 5}`

		expected := docker.HostConfig{}

		state := lua.NewState()
		So(lua.DoString(state, source), ShouldBeNil)
		state.Global("x")

		Convey("Then it is accepted and results in an unchaged result", func() {
			actual := ParseHostConfigFromLuaTable(state)
			So(actual, ShouldResemble, expected)
		})
	})
}
