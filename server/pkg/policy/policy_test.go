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

func TestParseTTLDuration(t *testing.T) {
	ttl, err := ParseTTLDuration("", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if ttl != time.Hour {
		t.Fatalf("expected fallback ttl, got %s", ttl)
	}
	ttl, err = ParseTTLDuration("7d", 0)
	if err != nil {
		t.Fatal(err)
	}
	if ttl != 7*24*time.Hour {
		t.Fatalf("expected 7d, got %s", ttl)
	}
	if _, err := ParseTTLDuration("bad", 0); err == nil {
		t.Fatal("expected invalid ttl to fail")
	}
}

func TestVisitLimitExceeded(t *testing.T) {
	if VisitLimitExceeded(0, 100) {
		t.Fatal("expected zero max visits to mean unlimited")
	}
	if !VisitLimitExceeded(3, 3) {
		t.Fatal("expected visit limit to be exceeded")
	}
}
