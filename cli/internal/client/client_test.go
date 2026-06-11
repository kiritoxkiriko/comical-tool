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
