package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

func GetEncryptPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 15000, 32, sha256.New)
	return base64.StdEncoding.EncodeToString(dk)
}
