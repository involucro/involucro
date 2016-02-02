package runtime

import "testing"

func TestRandomIdentifier(t *testing.T) {
	a := randomIdentifier()
	b := randomIdentifier()
	if a == b {
		t.Error("Same identifiers for a and b")
	}
}

func TestRandomIdentifierOfLength(t *testing.T) {
	s := randomIdentifierOfLength(64)
	if len(s) != 64 {
		t.Errorf("unexpected length: %v", len(s))
	}
}
