package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSha256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}