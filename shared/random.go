package shared

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomBytes(length uint) ([]byte, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func GenerateRandomString(length uint) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes)[:length], nil
}
