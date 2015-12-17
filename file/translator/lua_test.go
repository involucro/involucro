package translator

import (
	"github.com/Shopify/go-lua"
	"testing"
)

func TestLuaTableTraversal(t *testing.T) {
	state := lua.NewState()
	if err := lua.DoString(state, "tab = {a = 1, b = 2}"); err != nil {
		t.Fatal("Unable to run code", err)
	}
	state.Global("tab")
	state.PushNil()

	if !state.Next(-2) {
		t.Fatal("no next value")
	}

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
	if actual := lua.CheckString(state, -2); actual != key1 {
		t.Error("Unexpected key", actual)
	}
	if actual := lua.CheckString(state, -1); actual != val1 {
		t.Error("Unexpected value", actual)
	}
	state.Pop(1)

	if !state.Next(-2) {
		t.Fatal("No next value")
	}
	if actual := lua.CheckString(state, -2); actual != key2 {
		t.Error("Unexpected key", actual)
	}
	if actual := lua.CheckString(state, -1); actual != val2 {
		t.Error("Unexpected value", actual)
	}
	state.Pop(1)

	if state.Next(-2) {
		t.Fatal("Has next value")
	}
}
