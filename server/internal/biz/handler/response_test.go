package handler

import (
	"encoding/json"
	stdhttp "net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kiritoxkiriko/comical-tool/server/internal/biz/middleware"
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

func TestRequestIDReadsMiddlewareValue(t *testing.T) {
	c := &app.RequestContext{}
	c.Set(middleware.RequestIDKey, "req-123")
	if got := requestID(c); got != "req-123" {
		t.Fatalf("expected request id, got %q", got)
	}
}

func TestWriteResultWrapsBusinessData(t *testing.T) {
	c := app.NewContext(0)
	writeResult(c, map[string]string{"id": "abc"}, nil)

	var payload struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(c.Response.Body(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Data["id"] != "abc" {
		t.Fatalf("expected wrapped data, got %s", c.Response.Body())
	}
}

func TestWriteErrorWrapsCodeMessageAndRequestID(t *testing.T) {
	c := app.NewContext(0)
	c.Set(middleware.RequestIDKey, "req-123")
	writeError(c, apperror.New(apperror.CodeBadRequest, "invalid request"))

	var payload errorEnvelope
	if err := json.Unmarshal(c.Response.Body(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Error.Code != apperror.CodeBadRequest || payload.Error.Message != "invalid request" {
		t.Fatalf("unexpected error envelope: %+v", payload)
	}
	if payload.Error.RequestID != "req-123" {
		t.Fatalf("expected request id, got %q", payload.Error.RequestID)
	}
}
