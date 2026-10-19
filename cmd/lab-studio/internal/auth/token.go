package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomToken liefert 32 zufaellige Bytes als Hex-Text (Sitzungen, Mail-Links, CSRF).
func RandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
