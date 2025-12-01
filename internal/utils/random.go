package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string{
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seedRand  := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seedRand.Intn(len(charset))]
	}

	return string(b)
}


func CoalescePtr[T comparable](a,b *T) *T{
	if a != nil {
		return a 
	}

	return b
}