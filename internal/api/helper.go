package api

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)

	return base64.URLEncoding.EncodeToString(b), err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func extractID(r *http.Request, param string) (int, error) {
	paramString := chi.URLParam(r, param)

	return strconv.Atoi(paramString)
}
