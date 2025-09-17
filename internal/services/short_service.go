package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kiritoxkiriko/comical-tool/internal/config"
	"github.com/kiritoxkiriko/comical-tool/internal/database"
	"github.com/kiritoxkiriko/comical-tool/internal/models"
	redisClient "github.com/kiritoxkiriko/comical-tool/internal/redis"
	"gorm.io/gorm"
)

// ShortService handles short URL operations
type ShortService struct {
	config *config.Config
	db     *gorm.DB
	redis  *redis.Client
}

// NewShortService creates a new short service
func NewShortService(cfg *config.Config) *ShortService {
	return &ShortService{
		config: cfg,
		db:     database.GetDB(),
		redis:  redisClient.GetClient(),
	}
}

// CreateShort creates a new short URL
func (s *ShortService) CreateShort(ctx context.Context, req *models.CreateShortRequest) (*models.CreateShortResponse, error) {
	// Generate code if not provided
	code := req.CustomCode
	if code == nil || *code == "" {
		generatedCode, err := s.generateUniqueCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}
		code = &generatedCode
	}

	// Check if code already exists
	var existingShort models.Short
	err := s.db.Where("code = ?", *code).First(&existingShort).Error
	if err == nil {
		return nil, fmt.Errorf("code already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check code existence: %w", err)
	}

	// Set default expiry if not provided
	expiresAt := req.ExpiresAt
	if expiresAt == nil && s.config.ShortURL.DefaultExpiry > 0 {
		defaultExpiry := time.Now().Add(s.config.GetDefaultExpiry())
		expiresAt = &defaultExpiry
	}

	// Create short URL record
	short := &models.Short{
		Code:        *code,
		OriginalURL: req.OriginalURL,
		ExpiresAt:   expiresAt,
		MaxClicks:   req.MaxClicks,
		ClickCount:  0,
		IsActive:    true,
	}

	err = s.db.Create(short).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	// Cache the short URL
	cacheKey := fmt.Sprintf("short:%s", *code)
	cacheValue := short.OriginalURL
	cacheExpiry := 24 * time.Hour
	if expiresAt != nil {
		cacheExpiry = time.Until(*expiresAt)
	}

	err = redisClient.Set(ctx, cacheKey, cacheValue, cacheExpiry)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache short URL: %v\n", err)
	}

	// Build short URL
	shortURL := fmt.Sprintf("http://%s/%s", s.config.GetShortURLDomain(), *code)

	return &models.CreateShortResponse{
		ShortURL:    shortURL,
		OriginalURL: short.OriginalURL,
		CreatedAt:   short.CreatedAt,
		ExpiresAt:   short.ExpiresAt,
		MaxClicks:   short.MaxClicks,
	}, nil
}

// GetShort retrieves a short URL by code
func (s *ShortService) GetShort(ctx context.Context, code string) (*models.GetShortResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("short:%s", code)
	cachedURL, err := redisClient.Get(ctx, cacheKey)
	if err == nil && cachedURL != "" {
		// Get full record from database for complete response
		var short models.Short
		err = s.db.Where("code = ?", code).First(&short).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get short URL: %w", err)
		}

		shortURL := fmt.Sprintf("http://%s/%s", s.config.GetShortURLDomain(), code)
		return &models.GetShortResponse{
			ShortURL:    shortURL,
			OriginalURL: short.OriginalURL,
			CreatedAt:   short.CreatedAt,
			ExpiresAt:   short.ExpiresAt,
			MaxClicks:   short.MaxClicks,
			ClickCount:  short.ClickCount,
		}, nil
	}

	// Get from database
	var short models.Short
	err = s.db.Where("code = ? AND is_active = ?", code, true).First(&short).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("short URL not found")
		}
		return nil, fmt.Errorf("failed to get short URL: %w", err)
	}

	// Check if expired
	if short.ExpiresAt != nil && time.Now().After(*short.ExpiresAt) {
		return nil, fmt.Errorf("short URL has expired")
	}

	// Check click limit
	if short.MaxClicks != nil && short.ClickCount >= *short.MaxClicks {
		return nil, fmt.Errorf("short URL has reached maximum clicks")
	}

	// Cache the result
	cacheExpiry := 24 * time.Hour
	if short.ExpiresAt != nil {
		cacheExpiry = time.Until(*short.ExpiresAt)
	}
	redisClient.Set(ctx, cacheKey, short.OriginalURL, cacheExpiry)

	shortURL := fmt.Sprintf("http://%s/%s", s.config.GetShortURLDomain(), code)
	return &models.GetShortResponse{
		ShortURL:    shortURL,
		OriginalURL: short.OriginalURL,
		CreatedAt:   short.CreatedAt,
		ExpiresAt:   short.ExpiresAt,
		MaxClicks:   short.MaxClicks,
		ClickCount:  short.ClickCount,
	}, nil
}

// GetOriginalURL retrieves the original URL for redirection
func (s *ShortService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("short:%s", code)
	cachedURL, err := redisClient.Get(ctx, cacheKey)
	if err == nil && cachedURL != "" {
		return cachedURL, nil
	}

	// Get from database
	var short models.Short
	err = s.db.Where("code = ? AND is_active = ?", code, true).First(&short).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("short URL not found")
		}
		return "", fmt.Errorf("failed to get short URL: %w", err)
	}

	// Check if expired
	if short.ExpiresAt != nil && time.Now().After(*short.ExpiresAt) {
		return "", fmt.Errorf("short URL has expired")
	}

	// Check click limit
	if short.MaxClicks != nil && short.ClickCount >= *short.MaxClicks {
		return "", fmt.Errorf("short URL has reached maximum clicks")
	}

	// Cache the result
	cacheExpiry := 24 * time.Hour
	if short.ExpiresAt != nil {
		cacheExpiry = time.Until(*short.ExpiresAt)
	}
	redisClient.Set(ctx, cacheKey, short.OriginalURL, cacheExpiry)

	return short.OriginalURL, nil
}

// UpdateShort updates a short URL
func (s *ShortService) UpdateShort(ctx context.Context, code string, req *models.UpdateShortRequest) error {
	var short models.Short
	err := s.db.Where("code = ?", code).First(&short).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("short URL not found")
		}
		return fmt.Errorf("failed to get short URL: %w", err)
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.ExpiresAt != nil {
		updates["expires_at"] = *req.ExpiresAt
	}
	if req.MaxClicks != nil {
		updates["max_clicks"] = *req.MaxClicks
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	err = s.db.Model(&short).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update short URL: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("short:%s", code)
	redisClient.Del(ctx, cacheKey)

	return nil
}

// DeleteShort deletes a short URL
func (s *ShortService) DeleteShort(ctx context.Context, code string) error {
	var short models.Short
	err := s.db.Where("code = ?", code).First(&short).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("short URL not found")
		}
		return fmt.Errorf("failed to get short URL: %w", err)
	}

	err = s.db.Delete(&short).Error
	if err != nil {
		return fmt.Errorf("failed to delete short URL: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("short:%s", code)
	redisClient.Del(ctx, cacheKey)

	return nil
}

// generateUniqueCode generates a unique code
func (s *ShortService) generateUniqueCode() (string, error) {
	allowedChars := s.config.ShortURL.AllowedChars
	codeLength := s.config.ShortURL.CodeLength

	for i := 0; i < 10; i++ { // Try up to 10 times
		code := s.generateRandomCode(allowedChars, codeLength)

		// Check if code exists
		var count int64
		err := s.db.Model(&models.Short{}).Where("code = ?", code).Count(&count).Error
		if err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique code after 10 attempts")
}

// generateRandomCode generates a random code
func (s *ShortService) generateRandomCode(allowedChars string, length int) string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := range code {
		code[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return string(code)
}
