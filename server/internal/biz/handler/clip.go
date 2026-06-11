package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
)

type clipRequest struct {
	Content   string `json:"content"`
	Password  string `json:"password"`
	TTL       string `json:"ttl"`
	MaxVisits int    `json:"max_visits"`
	Link      bool   `json:"link"`
}

func (h *Handler) CreateClip(ctx context.Context, c *app.RequestContext) {
	var req clipRequest
	if !bindJSON(c, &req) {
		return
	}
	ttl, err := parseTTL(req.TTL, h.cfg.Modules.Clipboard.DefaultTTL)
	if err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "invalid ttl"))
		return
	}
	visits := req.MaxVisits
	if visits == 0 {
		visits = h.cfg.Modules.Clipboard.MaxVisits
	}
	item, err := h.svc.CreateClipboard(ctx, req.Content, req.Password, ttl, visits, req.Link)
	writeResult(c, item, err)
}

func (h *Handler) GetClip(ctx context.Context, c *app.RequestContext) {
	item, err := h.svc.GetClipboard(ctx, c.Param("id"), c.Query("password"))
	writeResult(c, item, err)
}

func (h *Handler) DeleteClip(ctx context.Context, c *app.RequestContext) {
	err := h.svc.DeleteClipboard(ctx, c.Param("id"))
	writeResult(c, utils.H{"deleted": err == nil}, err)
}
