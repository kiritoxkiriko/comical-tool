package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadWritesResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("expected bearer token, got %q", got)
		}
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	output := filepath.Join(t.TempDir(), "asset.bin")
	client := New(server.URL, "secret")
	if err := client.Download("/api/assets/id", output); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("expected downloaded body, got %q", string(data))
	}
}

func TestJSONUnwrapsDataEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"abc"}}`))
	}))
	defer server.Close()

	client := New(server.URL, "")
	data, err := client.JSON(http.MethodPost, "/api/test", map[string]string{"ok": "true"})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"id":"abc"}` {
		t.Fatalf("expected unwrapped data, got %s", data)
	}
}

func TestJSONFormatsEnvelopeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":"bad_request","message":"invalid","request_id":"req-1"}}`))
	}))
	defer server.Close()

	client := New(server.URL, "")
	_, err := client.JSON(http.MethodPost, "/api/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "400 Bad Request: invalid (request_id: req-1)" {
		t.Fatalf("unexpected error: %s", got)
	}
}
