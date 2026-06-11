package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
	"github.com/kiritoxkiriko/comical-tool/server/internal/repository"
	"github.com/kiritoxkiriko/comical-tool/server/internal/storage"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/policy"
)

type Service struct {
	cfg   config.Config
	repo  *repository.Store
	store storage.Store
}

type Upload struct {
	Name        string
	ContentType string
	Size        int64
	Body        io.Reader
	TTL         time.Duration
	Link        bool
}

func New(cfg config.Config, repo *repository.Store, store storage.Store) *Service {
	return &Service{cfg: cfg, repo: repo, store: store}
}

func (s *Service) CreateShortLink(ctx context.Context, targetURL, customSlug string, ttl time.Duration) (domain.ShortLink, error) {
	if _, err := url.ParseRequestURI(targetURL); err != nil {
		return domain.ShortLink{}, apperror.New(apperror.CodeBadRequest, "invalid target_url")
	}
	slug, err := s.slug(customSlug)
	if err != nil {
		return domain.ShortLink{}, err
	}
	id, err := policy.RandomID()
	if err != nil {
		return domain.ShortLink{}, err
	}
	link := domain.ShortLink{ID: id, Slug: slug, TargetURL: targetURL, ExpiresAt: policy.ExpiryFromDuration(ttl)}
	link = s.decorateShortLink(link)
	return link, s.repo.CreateShortLink(ctx, link)
}

func (s *Service) ResolveShortLink(ctx context.Context, slug string) (string, error) {
	link, err := s.repo.FindShortLink(ctx, slug)
	if errors.Is(err, repository.ErrNotFound) {
		return "", apperror.New(apperror.CodeNotFound, "short link not found")
	}
	if err != nil {
		return "", err
	}
	if link.RevokedAt != nil {
		return "", apperror.New(apperror.CodeRevoked, "short link revoked")
	}
	if policy.IsExpired(link.ExpiresAt) {
		return "", apperror.New(apperror.CodeExpired, "short link expired")
	}
	return link.TargetURL, nil
}

func (s *Service) RevokeShortLink(ctx context.Context, slug string) error {
	err := s.repo.RevokeShortLink(ctx, slug)
	if errors.Is(err, repository.ErrNotFound) {
		return apperror.New(apperror.CodeNotFound, "short link not found")
	}
	return err
}

func (s *Service) CreateClipboard(ctx context.Context, content, password string, ttl time.Duration, maxVisits int, link bool) (domain.ClipboardItem, error) {
	if strings.TrimSpace(content) == "" {
		return domain.ClipboardItem{}, apperror.New(apperror.CodeBadRequest, "content is required")
	}
	id, err := policy.RandomID()
	if err != nil {
		return domain.ClipboardItem{}, err
	}
	hash, err := policy.HashPassword(password)
	if err != nil {
		return domain.ClipboardItem{}, err
	}
	item := domain.ClipboardItem{ID: id, Content: content, PasswordHash: hash, MaxVisits: maxVisits, ExpiresAt: policy.ExpiryFromDuration(ttl)}
	if link {
		item.ShortSlug, err = s.linkTarget(ctx, s.publicURL("/api/clip/"+id), ttl)
	}
	if err != nil {
		return domain.ClipboardItem{}, err
	}
	return item, s.repo.CreateClipboard(ctx, item)
}

func (s *Service) GetClipboard(ctx context.Context, id, password string) (domain.ClipboardItem, error) {
	item, err := s.repo.FindClipboard(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return domain.ClipboardItem{}, apperror.New(apperror.CodeNotFound, "clipboard item not found")
	}
	if err != nil {
		return domain.ClipboardItem{}, err
	}
	if item.DeletedAt != nil || policy.IsExpired(item.ExpiresAt) {
		return domain.ClipboardItem{}, apperror.New(apperror.CodeExpired, "clipboard item expired")
	}
	if item.MaxVisits > 0 && item.VisitCount >= item.MaxVisits {
		return domain.ClipboardItem{}, apperror.New(apperror.CodeExpired, "clipboard item exhausted")
	}
	if !policy.CheckPassword(item.PasswordHash, password) {
		return domain.ClipboardItem{}, apperror.New(apperror.CodeForbidden, "invalid password")
	}
	if err := s.repo.IncrementClipboardVisit(ctx, id); err != nil {
		return domain.ClipboardItem{}, err
	}
	item.VisitCount++
	return item, nil
}

func (s *Service) DeleteClipboard(ctx context.Context, id string) error {
	err := s.repo.DeleteClipboard(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return apperror.New(apperror.CodeNotFound, "clipboard item not found")
	}
	return err
}

func (s *Service) UploadAsset(ctx context.Context, kind domain.ResourceType, up Upload) (domain.Asset, error) {
	if up.Size <= 0 {
		return domain.Asset{}, apperror.New(apperror.CodeBadRequest, "empty upload")
	}
	if limit := s.maxAssetBytes(kind); limit > 0 && up.Size > limit {
		return domain.Asset{}, apperror.New(apperror.CodeBadRequest, "upload exceeds max bytes")
	}
	asset, err := s.newAsset(ctx, kind, up)
	if err != nil {
		return domain.Asset{}, err
	}
	buf := &bytes.Buffer{}
	reader := io.TeeReader(io.LimitReader(up.Body, up.Size), buf)
	if err := s.store.Put(ctx, asset.ObjectKey, reader); err != nil {
		return domain.Asset{}, err
	}
	asset.Size = int64(buf.Len())
	return asset, s.repo.CreateAsset(ctx, asset)
}

func (s *Service) maxAssetBytes(kind domain.ResourceType) int64 {
	switch kind {
	case domain.ResourceImage:
		return s.cfg.Modules.ImageHosting.MaxBytes
	case domain.ResourceFile:
		return s.cfg.Modules.FileStash.MaxBytes
	default:
		return 0
	}
}

func (s *Service) ListAssets(ctx context.Context, kind domain.ResourceType) ([]domain.Asset, error) {
	return s.repo.ListAssets(ctx, kind)
}

func (s *Service) OpenAsset(ctx context.Context, id string) (domain.Asset, io.ReadCloser, error) {
	asset, err := s.repo.FindAsset(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return domain.Asset{}, nil, apperror.New(apperror.CodeNotFound, "asset not found")
	}
	if err != nil {
		return domain.Asset{}, nil, err
	}
	if asset.DeletedAt != nil || policy.IsExpired(asset.ExpiresAt) {
		return domain.Asset{}, nil, apperror.New(apperror.CodeExpired, "asset expired")
	}
	body, err := s.store.Open(ctx, asset.ObjectKey)
	return asset, body, err
}

func (s *Service) DeleteAsset(ctx context.Context, id string) error {
	asset, err := s.repo.FindAsset(ctx, id)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}
	if asset.ObjectKey != "" {
		_ = s.store.Delete(ctx, asset.ObjectKey)
	}
	err = s.repo.DeleteAsset(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return apperror.New(apperror.CodeNotFound, "asset not found")
	}
	return err
}

func (s *Service) CleanupExpired(ctx context.Context) (repository.CleanupResult, error) {
	return s.repo.CleanupExpired(ctx)
}

func (s *Service) newAsset(ctx context.Context, kind domain.ResourceType, up Upload) (domain.Asset, error) {
	id, err := policy.RandomID()
	if err != nil {
		return domain.Asset{}, err
	}
	asset := domain.Asset{
		ID: id, Kind: kind, Name: up.Name, ContentType: up.ContentType,
		Size: up.Size, ObjectKey: string(kind) + "/" + id,
		ExpiresAt: policy.ExpiryFromDuration(up.TTL),
	}
	if up.Link {
		asset.ShortSlug, err = s.linkTarget(ctx, s.publicURL("/api/assets/"+id), up.TTL)
	}
	return asset, err
}

func (s *Service) linkTarget(ctx context.Context, target string, ttl time.Duration) (string, error) {
	link, err := s.CreateShortLink(ctx, target, "", ttl)
	if err != nil {
		return "", err
	}
	return link.Slug, nil
}

func (s *Service) slug(custom string) (string, error) {
	if custom != "" {
		if !s.cfg.Modules.ShortLink.AllowCustomSlug || !policy.ValidateSlug(custom) {
			return "", apperror.New(apperror.CodeBadRequest, "invalid slug")
		}
		return custom, nil
	}
	return policy.RandomSlug()
}

func (s *Service) publicURL(path string) string {
	return strings.TrimRight(s.cfg.Server.PublicBaseURL, "/") + path
}

func (s *Service) decorateShortLink(link domain.ShortLink) domain.ShortLink {
	link.ShortURL = s.publicURL("/short/" + link.Slug)
	link.DomainURLs = map[string]string{}
	link.MappedURLs = map[string]string{}
	for host, mappedBase := range s.cfg.Modules.ShortLink.DomainMappings {
		link.DomainURLs[host] = "https://" + host + "/" + link.Slug
		link.MappedURLs[host] = strings.TrimRight(mappedBase, "/") + "/" + link.Slug
	}
	if len(link.DomainURLs) == 0 {
		link.DomainURLs = nil
		link.MappedURLs = nil
	}
	return link
}
