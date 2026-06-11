package main

import (
	"context"
	"flag"
	"log"

	"github.com/kiritoxkiriko/comical-tool/server/internal/biz/router"
	"github.com/kiritoxkiriko/comical-tool/server/internal/cache"
	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/internal/job"
	"github.com/kiritoxkiriko/comical-tool/server/internal/repository"
	"github.com/kiritoxkiriko/comical-tool/server/internal/service"
	"github.com/kiritoxkiriko/comical-tool/server/internal/storage"
)

func main() {
	configPath := flag.String("config", "", "path to config toml")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	cacheStore, err := cache.Open(ctx, cfg.Cache)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = cacheStore.Close()
	}()
	repo, err := repository.OpenSQLite(cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = repo.Close()
	}()
	if err := repo.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	store, err := storage.Open(ctx, cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}
	svc := service.New(cfg, repo, store)
	if cfg.Cleanup.Enabled {
		job.StartCleanup(ctx, cfg.Cleanup.Interval, svc, log.Default())
	}
	server := router.New(cfg, svc)
	server.Run()
}
