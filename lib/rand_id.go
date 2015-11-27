package lib

import (
	"crypto/rand"
	en "encoding/hex"
)

// RandomIdentifierOfLength returns a string
// composed of random integers (from `crypto/rand`)
// with exactly len characters.
func RandomIdentifierOfLength(len int) string {
	buflen := en.DecodedLen(len)
	buffer := make([]byte, buflen)
	rand.Read(buffer)

	return en.EncodeToString(buffer)
}

// RandomIdentifier is a shorthand for
// RandomIdentifierOfLength with a preset
// length
func RandomIdentifier() string {
	const LEN = 10
	return RandomIdentifierOfLength(LEN)
}
