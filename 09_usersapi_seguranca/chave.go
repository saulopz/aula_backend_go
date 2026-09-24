package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func gerarChave() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashChave(chave string) string {
	soma := sha256.Sum256([]byte(chave))
	return hex.EncodeToString(soma[:])
}
