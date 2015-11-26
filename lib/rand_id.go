package lib

import (
	"crypto/rand"
	en "encoding/hex"
)

func RandomIdentifier() string {
	const LEN = 10

	buffer := make([]byte, LEN)
	rand.Read(buffer)

	return en.EncodeToString(buffer)
}
