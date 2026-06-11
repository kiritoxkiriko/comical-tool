package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
)

type shortRequest struct {
	TargetURL  string `json:"target_url"`
	CustomSlug string `json:"custom_slug"`
	TTL        string `json:"ttl"`
}

func (h *Handler) CreateShort(ctx context.Context, c *app.RequestContext) {
	var req shortRequest
	if !bindJSON(c, &req) {
		return
	}
	ttl, err := parseTTL(req.TTL, h.cfg.Modules.ShortLink.DefaultTTL)
	if err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "invalid ttl"))
		return
	}
	link, err := h.svc.CreateShortLink(ctx, req.TargetURL, req.CustomSlug, ttl)
	writeResult(c, link, err)
}

func (h *Handler) RevokeShort(ctx context.Context, c *app.RequestContext) {
	err := h.svc.RevokeShortLink(ctx, c.Param("slug"))
	writeResult(c, utils.H{"revoked": err == nil}, err)
}

func (h *Handler) Redirect(ctx context.Context, c *app.RequestContext) {
	target, err := h.svc.ResolveShortLink(ctx, c.Param("slug"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.Response.Header.Set("Location", target)
	c.SetStatusCode(consts.StatusFound)
}
