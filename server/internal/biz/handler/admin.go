package handler

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
)

func (h *Handler) Cleanup(ctx context.Context, c *app.RequestContext) {
	if !adminAuthorized(h.cfg, string(c.Request.Header.Peek("Authorization"))) {
		writeError(c, apperror.New(apperror.CodeUnauthorized, "invalid admin token"))
		return
	}
	result, err := h.svc.CleanupExpired(ctx)
	writeResult(c, result, err)
}

func adminAuthorized(cfg config.Config, header string) bool {
	token := strings.TrimSpace(cfg.Security.AdminToken)
	if token == "" {
		return false
	}
	return strings.TrimSpace(header) == "Bearer "+token
}
