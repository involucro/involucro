package translator

import (
	"fmt"
	"github.com/Shopify/go-lua"
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
	l := lua.NewState()
	lua.DoString(l, `x = {"asd", "dsa", "x"}`)
	l.Global("x")
	m := checkStringSet(l, -1)
	if _, ok := m["asd"]; !ok {
		t.Error("asd not present")
	}
	if _, ok := m["dsa"]; !ok {
		t.Error("dsa not present")
	}
	if _, ok := m["x"]; !ok {
		t.Error("x not present")
	}
}

func TestCheckStringSetWithInteger(t *testing.T) {
	l := lua.NewState()

	lua.DoString(l, `x = 5`)
	l.Global("x")

	defer func() {
		if x := recover(); x == nil {
			t.Fatal("Didn't  panic")
		}
	}()
	checkStringSet(l, -1)
}
