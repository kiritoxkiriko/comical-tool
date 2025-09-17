# Product Requirements Document (PRD)
## Short URL Service

### 1. Executive Summary

The Short URL service is a core component of the comical-tool platform, providing URL shortening capabilities with advanced features like custom aliases, expiration controls, click tracking, and analytics. This service will enable users to create memorable, trackable short links while maintaining full control over their URL lifecycle.

### 2. Product Overview

**Product Name:** Short URL Service  
**Version:** 1.0  
**Target Users:** General web users, marketers, developers, content creators  
**Platform:** Backend API service (Golang + Hertz framework)

### 3. Business Objectives

- Provide a reliable, fast URL shortening service
- Enable users to create memorable, branded short links
- Offer comprehensive analytics for link performance tracking
- Support various expiration and access control mechanisms
- Build foundation for future comical-tool features

### 4. User Stories

#### Primary Users
- **Content Creator:** "I want to create short links for my social media posts with custom aliases that reflect my brand"
- **Marketer:** "I need to track click analytics and set expiration dates for campaign links"
- **Developer:** "I want an API to programmatically create and manage short URLs for my applications"

#### User Stories
1. As a user, I want to shorten a long URL so that I can share it easily
2. As a user, I want to create a custom alias for my short URL so that it's memorable and branded
3. As a user, I want to set an expiration time for my short URL so that it automatically becomes invalid
4. As a user, I want to limit the number of clicks before my URL expires so that I can control access
5. As a user, I want to view analytics for my short URLs so that I can track their performance
6. As a user, I want to configure URL generation options so that I can customize the service to my needs

### 5. Functional Requirements

#### 5.1 Core Features

**5.1.1 URL Shortening**
- Accept long URLs (up to 2048 characters)
- Generate short URLs with configurable length (default: 6-8 characters)
- Support custom character sets (alphanumeric, alphanumeric + symbols)
- Ensure URL uniqueness and collision handling
- Validate input URLs for proper format and accessibility

**5.1.2 Custom Code**
- Allow users to specify custom short codes
- Auto-generate codes if not provided
- Validate code format and availability
- Support case-insensitive codes
- Prevent reserved words and inappropriate content
- Enforce length limits (3-50 characters)

**5.1.3 Expiration Management**
- Time-based expiration (absolute date/time)
- Click-based expiration (expire after N clicks)
- Support for both expiration types simultaneously
- Automatic cleanup of expired URLs
- Graceful handling of expired URL access

**5.1.4 Analytics & Tracking**
- Click count tracking
- Referrer information capture
- User agent tracking
- Geographic data (country/city level)
- Timestamp of each click
- Click frequency over time
- Top referrers and user agents

**5.1.5 Configuration Options**
- Configurable URL length (3-20 characters)
- Custom character sets (alphanumeric, alphanumeric + symbols, custom)
- Default expiration settings
- Rate limiting configuration
- Analytics retention period
- Configurable short URL domain in config file (e.g., `short.example.com`, `go.mybrand.com`)
- Support for multiple domains with fallback to default
- Domain validation and SSL certificate management

#### 5.2 API Endpoints

**5.2.1 URL Management**
```
POST /api/v1/urls
- Create new short URL
- Request body: { original_url, custom_code?, expires_at?, max_clicks?, config? }
- Response: { short_url, original_url, created_at, expires_at?, max_clicks? }

GET /api/v1/urls/{code}
- Get URL details
- Response: { short_url, original_url, created_at, expires_at?, max_clicks?, click_count, analytics? }

PUT /api/v1/urls/{code}
- Update URL settings
- Request body: { expires_at?, max_clicks?, active? }
- Response: { success, message }

DELETE /api/v1/urls/{code}
- Delete short URL
- Response: { success, message }
```

**5.2.2 Analytics**
```
GET /api/v1/urls/{code}/analytics
- Get detailed analytics
- Query params: start_date?, end_date?, group_by?
- Response: { total_clicks, daily_clicks[], referrers[], countries[], user_agents[] }

GET /api/v1/urls/{code}/clicks
- Get click history
- Query params: page?, limit?, start_date?, end_date?
- Response: { clicks[], pagination }
```

**5.2.3 URL Redirection**
```
GET /{code}
- Redirect to original URL
- Handle expired/invalid URLs
- Track click analytics
- Response: 302 redirect or 404/410 error
```

#### 5.3 Data Models

**5.3.1 Short Entity**
```go
type Short struct {
    ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Code        string    `json:"code" gorm:"uniqueIndex"`
    OriginalURL string    `json:"original_url" gorm:"not null"`
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
    MaxClicks   *int      `json:"max_clicks,omitempty"`
    ClickCount  int       `json:"click_count" gorm:"default:0"`
    IsActive    bool      `json:"is_active" gorm:"default:true"`
    UserID      *string   `json:"user_id,omitempty"`
}
```

**5.3.2 Click Analytics Entity**
```go
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
```

### 6. Non-Functional Requirements

#### 6.1 Performance
- URL creation: < 100ms response time
- URL redirection: < 50ms response time
- Support 1000+ concurrent requests
- 99.9% uptime SLA

#### 6.2 Scalability
- Horizontal scaling capability
- Database sharding support
- Redis caching for hot URLs
- CDN integration for global performance

#### 6.3 Security
- Input validation and sanitization
- Rate limiting (100 requests/minute per IP)
- SQL injection prevention
- XSS protection
- HTTPS enforcement

#### 6.4 Reliability
- Automatic failover
- Data backup and recovery
- Graceful degradation
- Error monitoring and alerting

### 7. Technical Architecture

#### 7.1 Technology Stack
- **Backend:** Golang
- **Web Framework:** [Hertz](https://github.com/cloudwego/hertz) (CloudWeGo)
- **Database:** MySQL
- **Cache:** Redis
- **Containerization:** Docker

#### 7.2 System Components
- **Hertz HTTP Server:** Request routing and rate limiting
- **URL Service:** Core business logic
- **Analytics Service:** Click tracking and reporting
- **Configuration Service:** Domain and system configuration management
- **Cache Layer:** Redis for performance optimization
- **Database Layer:** MySQL for data persistence

#### 7.3 Data Flow
1. User creates short URL via API
2. System validates input and generates or uses provided code
3. URL data stored in MySQL (domain configured in config file)
4. Short URL returned to user (using configured domain)
5. When accessed via configured domain, system checks cache (Redis)
6. If not cached, queries database by code
7. Redirects to original URL
8. Tracks click analytics
9. Updates click count

### 8. Implementation Phases

#### Phase 1: Core Functionality (MVP)
- Basic URL shortening
- Custom code support (set or auto-generate)
- Time-based expiration
- Click counting
- Domain configuration in config file
- REST API endpoints

#### Phase 2: Advanced Features
- Click-based expiration
- Detailed analytics
- Configuration options
- Rate limiting
- Error handling improvements

#### Phase 3: Optimization & Scale
- Performance optimization
- Caching implementation
- Monitoring and alerting
- Load testing
- Documentation

### 9. Success Metrics

#### 9.1 Technical Metrics
- API response time < 100ms (95th percentile)
- System uptime > 99.9%
- Error rate < 0.1%
- Database query performance < 10ms

#### 9.2 Business Metrics
- URLs created per day
- Click-through rates
- User retention
- API adoption rate

### 10. Risk Assessment

#### 10.1 Technical Risks
- **Database performance:** Mitigation through indexing and caching
- **URL collision:** Mitigation through collision detection and retry logic
- **Scalability:** Mitigation through horizontal scaling and load balancing

#### 10.2 Security Risks
- **Malicious URLs:** Mitigation through URL validation and blacklisting
- **DDoS attacks:** Mitigation through rate limiting and CDN
- **Data breaches:** Mitigation through encryption and access controls

### 11. Future Enhancements

- User authentication and account management
- Bulk URL operations
- QR code generation
- Multiple domain support per instance
- Advanced analytics dashboard
- API rate limiting per user
- Webhook notifications
- URL preview generation
- Dynamic domain switching via API

---

**Document Version:** 1.0  
**Last Updated:** [Current Date]  
**Next Review:** [Date + 30 days]
