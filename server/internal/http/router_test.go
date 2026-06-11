package http

import (
	stdhttp "net/http"
	"testing"

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
