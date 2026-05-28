package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func New(prefix string) (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(b)), nil
}

func Must(prefix string) string {
	v, err := New(prefix)
	if err != nil {
		panic(err)
	}
	return v
}
