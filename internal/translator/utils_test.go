package translator

import (
	"testing"

	"github.com/Shopify/go-lua"
)

func TestCheckBoolean(t *testing.T) {
	l := lua.NewState()
	l.PushBoolean(true)
	if !checkBoolean(l, -1) {
		t.Error("expected true")
	}

	l.PushNumber(42)
	defer func() {
		err := recover()
		if err == nil {
			t.Error("expected panic")
		}
	}()
	checkBoolean(l, -1)
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

func TestCheckStringArray(t *testing.T) {
	l := lua.NewState()
	lua.DoString(l, `x = {"a=5", "b=6", "c=7"}`)
	l.Global("x")

	res := checkStringArray(l, -1)

	if len(res) != 3 {
		t.Fatal("Unexpected length of ", res)
	}
	if res[0] != "a=5" {
		t.Error("First element has wrong value", res[0])
	}
	if res[1] != "b=6" {
		t.Error("First element has wrong value", res[1])
	}
	if res[2] != "c=7" {
		t.Error("First element has wrong value", res[2])
	}
}
