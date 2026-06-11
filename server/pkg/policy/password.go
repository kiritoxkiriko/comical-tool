package policy

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	sum := sha256.Sum256(append(salt, []byte(password)...))
	return base64.RawStdEncoding.EncodeToString(salt) + "." +
		base64.RawStdEncoding.EncodeToString(sum[:]), nil
}

func CheckPassword(encoded string, password string) bool {
	if encoded == "" {
		return true
	}
	parts := strings.Split(encoded, ".")
	if len(parts) != 2 || password == "" {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	sum := sha256.Sum256(append(salt, []byte(password)...))
	return subtle.ConstantTimeCompare(sum[:], expected) == 1
}
