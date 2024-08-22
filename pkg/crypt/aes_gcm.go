package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func NewAESGCM(key Key) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher failed: %w", err)
	}

	return cipher.NewGCM(block)
}
