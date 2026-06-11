package domain

import "time"

const GuestUserID = "guest"

type ResourceType string

const (
	ResourceShortLink ResourceType = "short_link"
	ResourceImage     ResourceType = "image"
	ResourceClipboard ResourceType = "clipboard"
	ResourceFile      ResourceType = "file"
)

type ShortLink struct {
	ID         string            `json:"id"`
	Slug       string            `json:"slug"`
	TargetURL  string            `json:"target_url"`
	ShortURL   string            `json:"short_url,omitempty"`
	DomainURLs map[string]string `json:"domain_urls,omitempty"`
	MappedURLs map[string]string `json:"mapped_urls,omitempty"`
	ExpiresAt  *time.Time        `json:"expires_at,omitempty"`
	RevokedAt  *time.Time        `json:"revoked_at,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

type Asset struct {
	ID          string       `json:"id"`
	Kind        ResourceType `json:"kind"`
	Name        string       `json:"name"`
	ContentType string       `json:"content_type"`
	Size        int64        `json:"size"`
	ObjectKey   string       `json:"object_key"`
	ShortSlug   string       `json:"short_slug,omitempty"`
	ExpiresAt   *time.Time   `json:"expires_at,omitempty"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}

type ClipboardItem struct {
	ID           string     `json:"id"`
	Content      string     `json:"content,omitempty"`
	ShortSlug    string     `json:"short_slug,omitempty"`
	MaxVisits    int        `json:"max_visits"`
	VisitCount   int        `json:"visit_count"`
	PasswordHash string     `json:"-"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
