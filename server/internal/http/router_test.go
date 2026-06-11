package http

import (
	stdhttp "net/http"
	"testing"

	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
)

func TestStatusCodeMapsUnavailableResourcesToGone(t *testing.T) {
	cases := []apperror.Code{apperror.CodeExpired, apperror.CodeRevoked}
	for _, code := range cases {
		if got := statusCode(code); got != stdhttp.StatusGone {
			t.Fatalf("expected %s to map to 410, got %d", code, got)
		}
	}
}

func TestAdminAuthorizedRequiresBearerToken(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Security.AdminToken = "secret"
	if !adminAuthorized(cfg, "Bearer secret") {
		t.Fatal("expected valid bearer token to authorize")
	}
	if adminAuthorized(cfg, "secret") {
		t.Fatal("expected missing bearer scheme to fail")
	}
	if adminAuthorized(cfg, "Bearer wrong") {
		t.Fatal("expected wrong token to fail")
	}
	cfg.Security.AdminToken = ""
	if adminAuthorized(cfg, "Bearer secret") {
		t.Fatal("expected empty configured token to fail closed")
	}
}
