package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	stdhttp "net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/internal/service"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type Server struct {
	cfg config.Config
	svc *service.Service
}

func New(cfg config.Config, svc *service.Service) *Server {
	return &Server{cfg: cfg, svc: svc}
}

func (s *Server) Run() {
	h := s.engine()
	h.Use(cors)
	h.OPTIONS("/*path", options)
	h.GET("/healthz", s.health)
	h.GET("/short/:slug", s.redirect)
	h.GET("/:slug", s.redirect)
	api := h.Group("/api")
	api.POST("/short-links", s.createShort)
	api.POST("/short-links/:slug/revoke", s.revokeShort)
	api.POST("/clip", s.createClip)
	api.GET("/clip/:id", s.getClip)
	api.DELETE("/clip/:id", s.deleteClip)
	api.GET("/assets/:id", s.getAsset)
	api.POST("/images", s.uploadImage)
	api.GET("/images", s.listImages)
	api.DELETE("/images/:id", s.deleteAsset)
	api.POST("/files", s.uploadFile)
	api.GET("/files", s.listFiles)
	api.DELETE("/files/:id", s.deleteAsset)
	api.POST("/admin/cleanup", s.cleanup)
	h.Spin()
}

func (s *Server) engine() *server.Hertz {
	if s.cfg.Server.MaxBodyBytes <= 0 {
		return server.Default(server.WithHostPorts(s.cfg.Server.Addr))
	}
	return server.Default(
		server.WithHostPorts(s.cfg.Server.Addr),
		server.WithMaxRequestBodySize(int(s.cfg.Server.MaxBodyBytes)),
	)
}

func cors(ctx context.Context, c *app.RequestContext) {
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")
	c.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	c.Next(ctx)
}

func options(_ context.Context, c *app.RequestContext) {
	c.SetStatusCode(consts.StatusNoContent)
}

func (s *Server) health(_ context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"ok": true})
}

func (s *Server) createShort(ctx context.Context, c *app.RequestContext) {
	var req shortRequest
	if !bindJSON(c, &req) {
		return
	}
	ttl, err := parseTTL(req.TTL, s.cfg.Modules.ShortLink.DefaultTTL)
	if err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "invalid ttl"))
		return
	}
	link, err := s.svc.CreateShortLink(ctx, req.TargetURL, req.CustomSlug, ttl)
	writeResult(c, link, err)
}

func (s *Server) revokeShort(ctx context.Context, c *app.RequestContext) {
	err := s.svc.RevokeShortLink(ctx, c.Param("slug"))
	writeResult(c, utils.H{"revoked": err == nil}, err)
}

func (s *Server) redirect(ctx context.Context, c *app.RequestContext) {
	target, err := s.svc.ResolveShortLink(ctx, c.Param("slug"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.Response.Header.Set("Location", target)
	c.SetStatusCode(consts.StatusFound)
}

func (s *Server) createClip(ctx context.Context, c *app.RequestContext) {
	var req clipRequest
	if !bindJSON(c, &req) {
		return
	}
	ttl, err := parseTTL(req.TTL, s.cfg.Modules.Clipboard.DefaultTTL)
	if err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "invalid ttl"))
		return
	}
	visits := req.MaxVisits
	if visits == 0 {
		visits = s.cfg.Modules.Clipboard.MaxVisits
	}
	item, err := s.svc.CreateClipboard(ctx, req.Content, req.Password, ttl, visits, req.Link)
	writeResult(c, item, err)
}

func (s *Server) getClip(ctx context.Context, c *app.RequestContext) {
	item, err := s.svc.GetClipboard(ctx, c.Param("id"), c.Query("password"))
	writeResult(c, item, err)
}

func (s *Server) deleteClip(ctx context.Context, c *app.RequestContext) {
	err := s.svc.DeleteClipboard(ctx, c.Param("id"))
	writeResult(c, utils.H{"deleted": err == nil}, err)
}

func (s *Server) uploadImage(ctx context.Context, c *app.RequestContext) {
	s.uploadAsset(ctx, c, domain.ResourceImage, s.cfg.Modules.ImageHosting.DefaultTTL)
}

func (s *Server) uploadFile(ctx context.Context, c *app.RequestContext) {
	s.uploadAsset(ctx, c, domain.ResourceFile, s.cfg.Modules.FileStash.DefaultTTL)
}

func (s *Server) uploadAsset(ctx context.Context, c *app.RequestContext, kind domain.ResourceType, defaultTTL time.Duration) {
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
	asset, err := s.svc.UploadAsset(ctx, kind, up)
	writeResult(c, asset, err)
}

func (s *Server) listImages(ctx context.Context, c *app.RequestContext) {
	assets, err := s.svc.ListAssets(ctx, domain.ResourceImage)
	writeResult(c, assets, err)
}

func (s *Server) listFiles(ctx context.Context, c *app.RequestContext) {
	assets, err := s.svc.ListAssets(ctx, domain.ResourceFile)
	writeResult(c, assets, err)
}

func (s *Server) getAsset(ctx context.Context, c *app.RequestContext) {
	asset, body, err := s.svc.OpenAsset(ctx, c.Param("id"))
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

func (s *Server) deleteAsset(ctx context.Context, c *app.RequestContext) {
	err := s.svc.DeleteAsset(ctx, c.Param("id"))
	writeResult(c, utils.H{"deleted": err == nil}, err)
}

func (s *Server) cleanup(ctx context.Context, c *app.RequestContext) {
	result, err := s.svc.CleanupExpired(ctx)
	writeResult(c, result, err)
}

type shortRequest struct {
	TargetURL  string `json:"target_url"`
	CustomSlug string `json:"custom_slug"`
	TTL        string `json:"ttl"`
}

type clipRequest struct {
	Content   string `json:"content"`
	Password  string `json:"password"`
	TTL       string `json:"ttl"`
	MaxVisits int    `json:"max_visits"`
	Link      bool   `json:"link"`
}

func bindJSON(c *app.RequestContext, out any) bool {
	if err := json.Unmarshal(c.Request.Body(), out); err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "invalid json"))
		return false
	}
	return true
}

func parseTTL(value string, fallback time.Duration) (time.Duration, error) {
	if value == "" {
		return fallback, nil
	}
	return time.ParseDuration(value)
}

func writeResult(c *app.RequestContext, value any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, value)
}

func writeError(c *app.RequestContext, err error) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		appErr = apperror.New(apperror.CodeInternal, err.Error())
	}
	c.JSON(statusCode(appErr.Code), utils.H{"error": appErr.Code, "message": appErr.Message})
}

func statusCode(code apperror.Code) int {
	switch code {
	case apperror.CodeBadRequest:
		return stdhttp.StatusBadRequest
	case apperror.CodeNotFound:
		return stdhttp.StatusNotFound
	case apperror.CodeForbidden, apperror.CodeUnauthorized:
		return stdhttp.StatusForbidden
	case apperror.CodeConflict:
		return stdhttp.StatusConflict
	case apperror.CodeExpired, apperror.CodeRevoked:
		return stdhttp.StatusGone
	default:
		return stdhttp.StatusInternalServerError
	}
}
