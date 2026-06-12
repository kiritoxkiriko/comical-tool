package handler

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/kiritoxkiriko/comical-tool/server/internal/service"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

func (h *Handler) UploadImage(ctx context.Context, c *app.RequestContext) {
	h.uploadAsset(ctx, c, domain.ResourceImage, h.cfg.Modules.ImageHosting.DefaultTTL)
}

func (h *Handler) UploadFile(ctx context.Context, c *app.RequestContext) {
	h.uploadAsset(ctx, c, domain.ResourceFile, h.cfg.Modules.FileStash.DefaultTTL)
}

func (h *Handler) ListImages(ctx context.Context, c *app.RequestContext) {
	assets, err := h.svc.ListAssets(ctx, domain.ResourceImage)
	writeResult(c, assets, err)
}

func (h *Handler) ListFiles(ctx context.Context, c *app.RequestContext) {
	assets, err := h.svc.ListAssets(ctx, domain.ResourceFile)
	writeResult(c, assets, err)
}

func (h *Handler) GetAsset(ctx context.Context, c *app.RequestContext) {
	asset, body, err := h.svc.OpenAsset(ctx, c.Param("id"), c.Query("password"))
	if err != nil {
		writeError(c, err)
		return
	}
	defer func() {
		_ = body.Close()
	}()
	data, err := io.ReadAll(body)
	if err != nil {
		writeError(c, err)
		return
	}
	c.Data(consts.StatusOK, asset.ContentType, data)
}

func (h *Handler) DeleteAsset(ctx context.Context, c *app.RequestContext) {
	err := h.svc.DeleteAsset(ctx, c.Param("id"))
	writeResult(c, utils.H{"deleted": err == nil}, err)
}

func (h *Handler) uploadAsset(
	ctx context.Context,
	c *app.RequestContext,
	kind domain.ResourceType,
	defaultTTL time.Duration,
) {
	file, err := c.FormFile("file")
	if err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "file is required"))
		return
	}
	body, err := file.Open()
	if err != nil {
		writeError(c, err)
		return
	}
	defer func() {
		_ = body.Close()
	}()
	ttl, err := parseTTL(c.PostForm("ttl"), defaultTTL)
	if err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "invalid ttl"))
		return
	}
	up := service.Upload{
		Name: file.Filename, ContentType: file.Header.Get("Content-Type"),
		Size: file.Size, Body: body, TTL: ttl, Link: c.PostForm("link") == "true",
	}
	if kind == domain.ResourceFile {
		up.Password = c.PostForm("password")
		visits, err := parseOptionalInt(c.PostForm("max_visits"))
		if err != nil {
			writeError(c, apperror.New(apperror.CodeBadRequest, "invalid max_visits"))
			return
		}
		up.MaxVisits = visits
	}
	asset, err := h.svc.UploadAsset(ctx, kind, up)
	writeResult(c, asset, err)
}

func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if parsed < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return parsed, nil
}
