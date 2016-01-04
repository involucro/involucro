package utils

import (
	"fmt"
)

func ExampleRandomIdentifier_AreDifferent() {
	a := RandomIdentifier()
	b := RandomIdentifier()
	fmt.Println(a == b)
	// Output: false
}

func ExampleRandomIdentifierOfLength() {
	s := RandomIdentifierOfLength(64)
	fmt.Println(len(s))
	// Output: 64
}
