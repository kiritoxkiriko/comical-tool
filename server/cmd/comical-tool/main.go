package main

import (
	"context"
	"flag"
	"log"

	"github.com/kiritoxkiriko/comical-tool/server/internal/biz/router"
	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/internal/job"
	"github.com/kiritoxkiriko/comical-tool/server/internal/repository"
	"github.com/kiritoxkiriko/comical-tool/server/internal/service"
	"github.com/kiritoxkiriko/comical-tool/server/internal/storage"
)

func main() {
	configPath := flag.String("config", "", "path to config toml")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	repo, err := repository.OpenSQLite(cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = repo.Close()
	}()
	if err := repo.Migrate(context.Background()); err != nil {
		log.Fatal(err)
	}
	store := storage.NewLocal(cfg.Storage.LocalDir)
	svc := service.New(cfg, repo, store)
	if cfg.Cleanup.Enabled {
		job.StartCleanup(context.Background(), cfg.Cleanup.Interval, svc, log.Default())
	}
	server := router.New(cfg, svc)
	server.Run()
}
