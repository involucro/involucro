package translator

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseDockerConfig(t *testing.T) {
	Convey("Given an empty Lua state", t, func() {
		state := lua.NewState()

		run := func(s string) {
			So(lua.DoString(state, s), ShouldBeNil)
		}

		thenParseImageConfigPanics := func() {
			Convey("Then ParseImageConfigFromLuaTable panics", func() {
				So(func() { ParseImageConfigFromLuaTable(state) }, ShouldPanic)
			})
		}
		Convey("When there is nothing on the stack", thenParseImageConfigPanics)
		Convey("When there is a number on the stack", thenParseImageConfigPanics)

		Convey("When I set x = {} and parse it via ParseImageConfigFromLuaTable", func() {
			run(`x = {}`)
			state.Global("x")
			conf := ParseImageConfigFromLuaTable(state)

			Convey("Then the parsed config contains no Cmd", func() {
				So(conf.Cmd, ShouldHaveLength, 0)
			})
		})

		Convey("When I set x = {User = 'blah'} and parse it via ParseImageConfigFromLuaTable", func() {
			run(`x = {User = 'blah'}`)
			state.Global("x")
			conf := ParseImageConfigFromLuaTable(state)

			Convey("Then the parsed config contains no Cmd", func() {
				So(conf.Cmd, ShouldHaveLength, 0)
			})
			Convey("Then the parsed config has the given User", func() {
				So(conf.User, ShouldResemble, "blah")
			})
		})
		Convey("When I set x = {Cmd = 5} and parse it via ParseImageConfigFromLuaTable", thenParseImageConfigPanics)

		Convey("When I set x = {User = 'blah', Cmd = ['asd', dsa']} and parse it via ParseImageConfigFromLuaTable", func() {
			run(`x = {User = 'blah', Cmd = {'asd', 'dsa'}}`)
			state.Global("x")
			conf := ParseImageConfigFromLuaTable(state)

			Convey("Then the parsed config contains the Cmd ['asd', 'dsa']", func() {
				So(conf.Cmd, ShouldResemble, []string{"asd", "dsa"})
			})
			Convey("Then the parsed config has the given User", func() {
				So(conf.User, ShouldResemble, "blah")
			})
		})
		Convey("When I set x = {Cmd = ['asd', dsa'], Entrypoint = ['ttt', '123']} and parse it via ParseImageConfigFromLuaTable", func() {
			run(`x = {Cmd = {'asd', 'dsa'}, Entrypoint = {'ttt', '123'}}`)
			state.Global("x")
			conf := ParseImageConfigFromLuaTable(state)

			Convey("Then the parsed config contains the Cmd ['asd', 'dsa']", func() {
				So(conf.Cmd, ShouldResemble, []string{"asd", "dsa"})
			})
			Convey("Then the parsed config contains the Entrypoint ['ttt', '123']", func() {
				So(conf.Entrypoint, ShouldResemble, []string{"ttt", "123"})
			})
		})

		Convey("When I set x = {Labels = {testlabel = 'asd'}} and parse it", func() {
			run(`x = {Labels = {testlabel = 'asd'}}`)
			state.Global("x")
			conf := ParseImageConfigFromLuaTable(state)

			Convey("Then the parsed config has the label testlabel", func() {
				So(conf.Labels["testlabel"], ShouldResemble, "asd")
			})
		})
		Convey("When I set x = {ExposedPorts = {'80/tcp', '22/udp'}} and parse it", func() {
			run(`x = {ExposedPorts = {'80/tcp', '22/udp'}}`)
			state.Global("x")
			conf := ParseImageConfigFromLuaTable(state)

			Convey("Then the parsed config has those ports exposed", func() {
				So(conf.ExposedPorts, ShouldContainKey, docker.Port("80/tcp"))
				So(conf.ExposedPorts, ShouldContainKey, docker.Port("22/udp"))
			})
		})
		Convey("When I set x = {Volumes = {'/data'}} and parse it", func() {
			run(`x = {Volumes = {'/data'}}`)
			state.Global("x")
			conf := ParseImageConfigFromLuaTable(state)

			Convey("Then the parsed config has that volume", func() {
				So(conf.Volumes, ShouldContainKey, "/data")
			})
		})
	})
}
