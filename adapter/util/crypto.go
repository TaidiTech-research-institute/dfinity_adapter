package util

import "crypto/sha256"

func Sha256(m []byte) []byte {
	hash := sha256.Sum256(m)
	return hash[:]
}
