package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/internal/service"
)

type Handler struct {
	cfg config.Config
	svc *service.Service
}

func New(cfg config.Config, svc *service.Service) *Handler {
	return &Handler{cfg: cfg, svc: svc}
}

func (h *Handler) Health(_ context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"ok": true})
}
