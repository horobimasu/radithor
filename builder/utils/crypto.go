package utils

import (
	"crypto/ecdh"
	"crypto/rand"
)

func GenerateE2eeKey() (*ecdh.PrivateKey, error) {
	key, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return key, nil
}
