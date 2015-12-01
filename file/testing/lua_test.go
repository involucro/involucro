package testing

import (
	"github.com/Shopify/go-lua"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLuaTableTraversal(t *testing.T) {
	Convey("Given I have a Lua Context", t, func() {
		state := lua.NewState()
		Convey("When I put some keys into a table", func() {
			err := lua.DoString(state, "tab = {a = 1, b = 2}")
			So(err, ShouldBeNil)
			Convey("Then I can iterate over the keys", func() {
				state.Global("tab")
				state.PushNil()

				So(state.Next(-2), ShouldBeTrue)

				key1 := lua.CheckString(state, -2)
				var val1, key2, val2 string

				if key1 == "a" {
					val1 = "1"
					key2 = "b"
					val2 = "2"
				} else {
					val1 = "2"
					key2 = "a"
					val2 = "1"
				}
				So(lua.CheckString(state, -2), ShouldResemble, key1)
				So(lua.CheckString(state, -1), ShouldResemble, val1)
				state.Pop(1)

				So(state.Next(-2), ShouldBeTrue)
				So(lua.CheckString(state, -2), ShouldResemble, key2)
				So(lua.CheckString(state, -1), ShouldResemble, val2)
				state.Pop(1)

				So(state.Next(-2), ShouldBeFalse)
			})
		})
	})
}
