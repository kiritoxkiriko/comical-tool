package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/internal/repository"
	"github.com/kiritoxkiriko/comical-tool/server/internal/storage"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

func TestUploadAssetRejectsImageOverModuleLimit(t *testing.T) {
	svc := newTestService(t)
	svc.cfg.Modules.ImageHosting.MaxBytes = 5
	up := Upload{
		Name: "large.png", ContentType: "image/png", Size: 6,
		Body: strings.NewReader("123456"),
	}
	_, err := svc.UploadAsset(context.Background(), domain.ResourceImage, up)
	if !hasAppCode(err, apperror.CodeBadRequest) {
		t.Fatalf("expected bad_request, got %v", err)
	}
}

func TestUploadAssetAllowsFileWithinModuleLimit(t *testing.T) {
	svc := newTestService(t)
	svc.cfg.Modules.FileStash.MaxBytes = 6
	up := Upload{
		Name: "small.txt", ContentType: "text/plain", Size: 5,
		Body: strings.NewReader("12345"),
	}
	asset, err := svc.UploadAsset(context.Background(), domain.ResourceFile, up)
	if err != nil {
		t.Fatal(err)
	}
	if asset.Size != 5 {
		t.Fatalf("expected stored size 5, got %d", asset.Size)
	}
}

func newTestService(t *testing.T) *Service {
	t.Helper()
	cfg, err := config.Load("")
	if err != nil {
		t.Fatal(err)
	}
	repo, err := repository.OpenSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	if err := repo.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	return New(cfg, repo, storage.NewLocal(t.TempDir()))
}

func hasAppCode(err error, code apperror.Code) bool {
	var appErr *apperror.Error
	return errors.As(err, &appErr) && appErr.Code == code
}
