package service

import (
	"context"
	"errors"
	"io"
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

func TestUploadAssetRejectsNonImageMimeForImages(t *testing.T) {
	svc := newTestService(t)
	up := Upload{
		Name: "not-image.txt", ContentType: "text/plain", Size: 5,
		Body: strings.NewReader("hello"),
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

func TestOpenFileAssetRequiresPasswordAndVisitLimit(t *testing.T) {
	svc := newTestService(t)
	up := Upload{
		Name: "secret.txt", ContentType: "text/plain", Size: 6,
		Body: strings.NewReader("secret"), Password: "open", MaxVisits: 1,
	}
	asset, err := svc.UploadAsset(context.Background(), domain.ResourceFile, up)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.OpenAsset(context.Background(), asset.ID, "wrong"); !hasAppCode(err, apperror.CodeForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	opened, body, err := svc.OpenAsset(context.Background(), asset.ID, "open")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = body.Close()
	}()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "secret" {
		t.Fatalf("expected stored body, got %q", data)
	}
	if opened.VisitCount != 1 {
		t.Fatalf("expected visit count 1, got %d", opened.VisitCount)
	}
	if _, _, err := svc.OpenAsset(context.Background(), asset.ID, "open"); !hasAppCode(err, apperror.CodeExpired) {
		t.Fatalf("expected exhausted asset, got %v", err)
	}
}

func TestResolveShortLinkRecordsAccessEvent(t *testing.T) {
	svc := newTestService(t)
	link, err := svc.CreateShortLink(context.Background(), "https://example.com", "tracked", 0)
	if err != nil {
		t.Fatal(err)
	}
	target, err := svc.ResolveShortLink(context.Background(), "tracked")
	if err != nil {
		t.Fatal(err)
	}
	if target != "https://example.com" {
		t.Fatalf("expected target URL, got %q", target)
	}
	events, err := svc.repo.ListAccessEvents(context.Background(), domain.ResourceShortLink, link.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 access event, got %d", len(events))
	}
	if events[0].Action != "redirect" {
		t.Fatalf("expected redirect event, got %+v", events[0])
	}
}

func TestCreateShortLinkRejectsDuplicateActiveSlug(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.CreateShortLink(context.Background(), "https://example.com/a", "dup", 0); err != nil {
		t.Fatal(err)
	}
	_, err := svc.CreateShortLink(context.Background(), "https://example.com/b", "dup", 0)
	if !hasAppCode(err, apperror.CodeConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestCreateClipboardWithLinkRecordsResourceLink(t *testing.T) {
	svc := newTestService(t)
	item, err := svc.CreateClipboard(context.Background(), "hello", "", 0, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if item.ShortSlug == "" {
		t.Fatal("expected short slug")
	}
	links, err := svc.repo.ListResourceLinks(context.Background(), domain.ResourceClipboard, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 resource link, got %d", len(links))
	}
	if links[0].ResourceID != item.ID || links[0].ResourceType != domain.ResourceClipboard {
		t.Fatalf("unexpected resource link: %+v", links[0])
	}
}

func TestUploadAssetWithLinkRecordsResourceLink(t *testing.T) {
	svc := newTestService(t)
	asset, err := svc.UploadAsset(context.Background(), domain.ResourceFile, Upload{
		Name: "linked.txt", ContentType: "text/plain", Size: 6,
		Body: strings.NewReader("linked"), Link: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if asset.ShortSlug == "" {
		t.Fatal("expected short slug")
	}
	links, err := svc.repo.ListResourceLinks(context.Background(), domain.ResourceFile, asset.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 resource link, got %d", len(links))
	}
	if links[0].ResourceID != asset.ID || links[0].ResourceType != domain.ResourceFile {
		t.Fatalf("unexpected resource link: %+v", links[0])
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
