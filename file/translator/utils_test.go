package translator

import (
	"fmt"
	"github.com/Shopify/go-lua"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func ExampleCheckBoolean() {
	l := lua.NewState()
	l.PushBoolean(true)
	fmt.Printf("%t\n", checkBoolean(l, -1))
	// Output: true
}

func ExampleCheckBoolean_Failing() {
	l := lua.NewState()
	l.PushNumber(42)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("error occurred")
		}
	}()
	checkBoolean(l, -1)
	// Output: error occurred
}

func TestCheckStringSet(t *testing.T) {
	Convey("When I parse a string table as string set", t, func() {
		l := lua.NewState()
		lua.DoString(l, `x = {"asd", "dsa", "x"}`)
		l.Global("x")
		m := checkStringSet(l, -1)
		Convey("Then the keys are present", func() {
			So(m, ShouldContainKey, "asd")
			So(m, ShouldContainKey, "dsa")
			So(m, ShouldContainKey, "x")
		})
	})
	Convey("When I parse a number as string set", t, func() {
		l := lua.NewState()
		lua.DoString(l, `x = 5`)
		l.Global("x")
		Convey("Then it panics", func() {
			So(func() { checkStringSet(l, -1) }, ShouldPanic)
		})
	})
}
