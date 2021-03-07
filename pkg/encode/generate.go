package encode

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate a random base64 url encoded string
func GenerateBase64ID(size int, prefix string) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
