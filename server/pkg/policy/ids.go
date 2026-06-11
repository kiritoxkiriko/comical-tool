package policy

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
)

var slugPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`)

func RandomID() (string, error) {
	return randomHex(16)
}

func RandomSlug() (string, error) {
	value, err := randomHex(4)
	if err != nil {
		return "", err
	}
	return value, nil
}

func ValidateSlug(slug string) bool {
	return slugPattern.MatchString(slug)
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
