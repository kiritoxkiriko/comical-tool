package policy

import (
	"testing"
	"time"
)

func TestPasswordHash(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, "secret") {
		t.Fatal("expected password to match")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("expected wrong password to fail")
	}
}

func TestSlugValidation(t *testing.T) {
	if !ValidateSlug("abc_123-test") {
		t.Fatal("expected slug to be valid")
	}
	if ValidateSlug("ab") {
		t.Fatal("expected short slug to fail")
	}
}

func TestExpiry(t *testing.T) {
	expires := ExpiryFromDuration(time.Millisecond)
	if expires == nil {
		t.Fatal("expected expiry")
	}
	time.Sleep(2 * time.Millisecond)
	if !IsExpired(expires) {
		t.Fatal("expected expired value")
	}
}
