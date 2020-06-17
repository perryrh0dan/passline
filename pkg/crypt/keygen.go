package crypt

import (
	"crypto/rand"
	"io"
)

func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}

	return string(key), nil
}
