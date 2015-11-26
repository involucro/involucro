package lib

import (
	"crypto/rand"
	en "encoding/hex"
)

func RandomIdentifierOfLength(len int) string {
	buffer := make([]byte, len)
	rand.Read(buffer)

	return en.EncodeToString(buffer)
}

func RandomIdentifier() string {
	const LEN = 10
	return RandomIdentifierOfLength(LEN)
}
