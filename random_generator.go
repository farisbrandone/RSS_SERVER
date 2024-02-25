package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, errors.New("probleme to random data")
	}

	hmac := hmac.New(sha256.New, b)
	hmac.Write([]byte(b))
	dataHmac := hmac.Sum(nil)
	log.Println((dataHmac))
	return dataHmac, nil
}
