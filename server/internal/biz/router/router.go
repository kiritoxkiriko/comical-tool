package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/kiritoxkiriko/comical-tool/server/internal/biz/handler"
	"github.com/kiritoxkiriko/comical-tool/server/internal/biz/middleware"
	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/internal/service"
)

type Server struct {
	cfg config.Config
	svc *service.Service
}

func New(cfg config.Config, svc *service.Service) *Server {
	return &Server{cfg: cfg, svc: svc}
}

func (s *Server) Run() {
	s.Engine().Spin()
}

func (s *Server) Engine() *server.Hertz {
	h := s.newEngine()
	routes := handler.New(s.cfg, s.svc)
	h.Use(middleware.RequestID, middleware.CORS)
	h.OPTIONS("/*path", middleware.Options)
	h.GET("/healthz", routes.Health)
	h.GET("/api/health", routes.Health)
	h.GET("/short/:slug", routes.Redirect)
	h.GET("/:slug", routes.Redirect)
	api := h.Group("/api")
	api.POST("/short-links", routes.CreateShort)
	api.POST("/short-links/:slug/revoke", routes.RevokeShort)
	api.POST("/clip", routes.CreateClip)
	api.GET("/clip/:id", routes.GetClip)
	api.DELETE("/clip/:id", routes.DeleteClip)
	api.GET("/assets/:id", routes.GetAsset)
	api.POST("/images", routes.UploadImage)
	api.GET("/images", routes.ListImages)
	api.DELETE("/images/:id", routes.DeleteAsset)
	api.POST("/files", routes.UploadFile)
	api.GET("/files", routes.ListFiles)
	api.DELETE("/files/:id", routes.DeleteAsset)
	api.POST("/admin/cleanup", routes.Cleanup)
	return h
}

func (s *Server) newEngine() *server.Hertz {
	if s.cfg.Server.MaxBodyBytes <= 0 {
		return server.Default(server.WithHostPorts(s.cfg.Server.Addr))
	}
	return server.Default(
		server.WithHostPorts(s.cfg.Server.Addr),
		server.WithMaxRequestBodySize(int(s.cfg.Server.MaxBodyBytes)),
	)
}
