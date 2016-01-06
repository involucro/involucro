package runtime

import (
	"fmt"
)

func ExampleRandomIdentifier_AreDifferent() {
	a := randomIdentifier()
	b := randomIdentifier()
	fmt.Println(a == b)
	// Output: false
}

func ExampleRandomIdentifierOfLength() {
	s := randomIdentifierOfLength(64)
	fmt.Println(len(s))
	// Output: 64
}
