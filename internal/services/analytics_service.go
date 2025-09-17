package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kiritoxkiriko/comical-tool/internal/database"
	"github.com/kiritoxkiriko/comical-tool/internal/models"
	"gorm.io/gorm"
)

// AnalyticsService handles analytics operations
type AnalyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		db: database.GetDB(),
	}
}

// TrackClick records a click event
func (s *AnalyticsService) TrackClick(ctx context.Context, code, ipAddress, userAgent, referrer, country, city string) error {
	// Create click analytics record
	click := &models.ClickAnalytics{
		ID:        uuid.New().String(),
		URLCode:   code,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Referrer:  referrer,
		Country:   country,
		City:      city,
		ClickedAt: time.Now(),
	}

	err := s.db.Create(click).Error
	if err != nil {
		return fmt.Errorf("failed to track click: %w", err)
	}

	// Update click count in short URL
	err = s.db.Model(&models.Short{}).
		Where("code = ?", code).
		Update("click_count", gorm.Expr("click_count + 1")).Error
	if err != nil {
		return fmt.Errorf("failed to update click count: %w", err)
	}

	return nil
}

// GetAnalytics retrieves analytics for a short URL
func (s *AnalyticsService) GetAnalytics(ctx context.Context, code string, startDate, endDate *time.Time) (*models.AnalyticsResponse, error) {
	// Check if short URL exists
	var short models.Short
	err := s.db.Where("code = ?", code).First(&short).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("short URL not found")
		}
		return nil, fmt.Errorf("failed to get short URL: %w", err)
	}

	// Build query
	query := s.db.Model(&models.ClickAnalytics{}).Where("url_code = ?", code)

	if startDate != nil {
		query = query.Where("clicked_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("clicked_at <= ?", *endDate)
	}

	// Get total clicks
	var totalClicks int64
	err = query.Count(&totalClicks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get daily clicks
	var dailyClicks []models.DailyClickData
	err = query.Select("DATE(clicked_at) as date, COUNT(*) as clicks").
		Group("DATE(clicked_at)").
		Order("date").
		Scan(&dailyClicks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get daily clicks: %w", err)
	}

	// Get referrers
	var referrers []models.ReferrerData
	err = query.Select("referrer, COUNT(*) as count").
		Where("referrer != ''").
		Group("referrer").
		Order("count DESC").
		Limit(10).
		Scan(&referrers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get referrers: %w", err)
	}

	// Get countries
	var countries []models.CountryData
	err = query.Select("country, COUNT(*) as count").
		Where("country != ''").
		Group("country").
		Order("count DESC").
		Limit(10).
		Scan(&countries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get countries: %w", err)
	}

	// Get user agents
	var userAgents []models.UserAgentData
	err = query.Select("user_agent, COUNT(*) as count").
		Where("user_agent != ''").
		Group("user_agent").
		Order("count DESC").
		Limit(10).
		Scan(&userAgents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user agents: %w", err)
	}

	return &models.AnalyticsResponse{
		TotalClicks: int(totalClicks),
		DailyClicks: dailyClicks,
		Referrers:   referrers,
		Countries:   countries,
		UserAgents:  userAgents,
	}, nil
}

// GetClickHistory retrieves click history for a short URL
func (s *AnalyticsService) GetClickHistory(ctx context.Context, code string, page, limit int, startDate, endDate *time.Time) (*models.ClickHistoryResponse, error) {
	// Check if short URL exists
	var short models.Short
	err := s.db.Where("code = ?", code).First(&short).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("short URL not found")
		}
		return nil, fmt.Errorf("failed to get short URL: %w", err)
	}

	// Build query
	query := s.db.Model(&models.ClickAnalytics{}).Where("url_code = ?", code)

	if startDate != nil {
		query = query.Where("clicked_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("clicked_at <= ?", *endDate)
	}

	// Get total count
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Calculate pagination
	offset := (page - 1) * limit
	pages := int((total + int64(limit) - 1) / int64(limit))

	// Get clicks
	var clicks []models.ClickData
	err = query.Select("ip_address, user_agent, referrer, country, city, clicked_at").
		Order("clicked_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&clicks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks: %w", err)
	}

	return &models.ClickHistoryResponse{
		Clicks: clicks,
		Pagination: models.Pagination{
			Page:  page,
			Limit: limit,
			Total: int(total),
			Pages: pages,
		},
	}, nil
}

// CleanupOldAnalytics removes old analytics data
func (s *AnalyticsService) CleanupOldAnalytics(ctx context.Context, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	err := s.db.Where("clicked_at < ?", cutoffDate).Delete(&models.ClickAnalytics{}).Error
	if err != nil {
		return fmt.Errorf("failed to cleanup old analytics: %w", err)
	}

	return nil
}
