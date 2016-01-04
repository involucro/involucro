package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

var (
	encoding = base64.RawURLEncoding
)

// RegisterEncodeableType has to be called for
// any type that is part of an encoded state.
// See EncodeState for encoding.
func RegisterEncodeableType(v interface{}) {
	gob.Register(v)
}

func EncodeState(steps []Step) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(steps)
	encodedState := buf.Bytes()

	return encoding.EncodeToString(encodedState)
}

func DecodeState(state string) ([]Step, error) {
	debased, err := encoding.DecodeString(state)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(debased)
	dec := gob.NewDecoder(buffer)

	var steps []Step
	err = dec.Decode(&steps)
	return steps, err
}
