package crypt

import "crypto/sha256"

func NewKey(secret string) Key {
	a := sha256.Sum256([]byte(secret))
	return a[:]
}
