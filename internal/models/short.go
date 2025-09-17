package models

import (
	"time"
)

// Short represents a short URL record
type Short struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string     `json:"code" gorm:"uniqueIndex;not null"`
	OriginalURL string     `json:"original_url" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int       `json:"max_clicks,omitempty"`
	ClickCount  int        `json:"click_count" gorm:"default:0"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	UserID      *string    `json:"user_id,omitempty"`
}

// ClickAnalytics represents click tracking data
type ClickAnalytics struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	URLCode   string    `json:"url_code" gorm:"index"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	ClickedAt time.Time `json:"clicked_at"`
}

// CreateShortRequest represents the request to create a short URL
type CreateShortRequest struct {
	OriginalURL string     `json:"original_url" binding:"required,url"`
	CustomCode  *string    `json:"custom_code,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int       `json:"max_clicks,omitempty"`
}

// CreateShortResponse represents the response when creating a short URL
type CreateShortResponse struct {
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int       `json:"max_clicks,omitempty"`
}

// GetShortResponse represents the response when getting short URL details
type GetShortResponse struct {
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int       `json:"max_clicks,omitempty"`
	ClickCount  int        `json:"click_count"`
}

// UpdateShortRequest represents the request to update a short URL
type UpdateShortRequest struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	MaxClicks *int       `json:"max_clicks,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}

// AnalyticsResponse represents analytics data
type AnalyticsResponse struct {
	TotalClicks int              `json:"total_clicks"`
	DailyClicks []DailyClickData `json:"daily_clicks"`
	Referrers   []ReferrerData   `json:"referrers"`
	Countries   []CountryData    `json:"countries"`
	UserAgents  []UserAgentData  `json:"user_agents"`
}

// DailyClickData represents daily click statistics
type DailyClickData struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

// ReferrerData represents referrer statistics
type ReferrerData struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

// CountryData represents country statistics
type CountryData struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

// UserAgentData represents user agent statistics
type UserAgentData struct {
	UserAgent string `json:"user_agent"`
	Count     int    `json:"count"`
}

// ClickHistoryResponse represents click history data
type ClickHistoryResponse struct {
	Clicks     []ClickData `json:"clicks"`
	Pagination Pagination  `json:"pagination"`
}

// ClickData represents individual click data
type ClickData struct {
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	ClickedAt time.Time `json:"clicked_at"`
}

// Pagination represents pagination information
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}
