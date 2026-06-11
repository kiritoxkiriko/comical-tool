package main

import (
	"context"
	"flag"
	"log"

	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	apihttp "github.com/kiritoxkiriko/comical-tool/server/internal/http"
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
	server := apihttp.New(cfg, service.New(cfg, repo, store))
	server.Run()
}
