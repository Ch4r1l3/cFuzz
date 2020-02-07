package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b)[:length], err
}
