namespace go comical_tool

// Short URL service definitions
struct ShortURL {
    1: required i32 id,
    2: required string code,
    3: required string original_url,
    4: required string created_at,
    5: optional string expires_at,
    6: optional i32 max_clicks,
    7: required i32 click_count,
    8: required bool is_active,
    9: optional string user_id,
}

struct CreateShortRequest {
    1: required string original_url (api.body="original_url"),
    2: optional string custom_code (api.body="custom_code"),
    3: optional string expires_at (api.body="expires_at"),
    4: optional i32 max_clicks (api.body="max_clicks"),
}

struct CreateShortResponse {
    1: required string short_url,
    2: required string original_url,
    3: required string created_at,
    4: optional string expires_at,
    5: optional i32 max_clicks,
}

struct GetShortResponse {
    1: required string short_url,
    2: required string original_url,
    3: required string created_at,
    4: optional string expires_at,
    5: optional i32 max_clicks,
    6: required i32 click_count,
}

struct UpdateShortRequest {
    1: required string code (api.path="code"),
    2: optional string expires_at (api.body="expires_at"),
    3: optional i32 max_clicks (api.body="max_clicks"),
    4: optional bool is_active (api.body="is_active"),
}

struct UpdateShortResponse {
    1: required bool success,
    2: required string message,
}

struct DeleteShortResponse {
    1: required bool success,
    2: required string message,
}

struct DailyClickData {
    1: required string date,
    2: required i32 clicks,
}

struct ReferrerData {
    1: required string referrer,
    2: required i32 count,
}

struct CountryData {
    1: required string country,
    2: required i32 count,
}

struct UserAgentData {
    1: required string user_agent,
    2: required i32 count,
}

struct AnalyticsResponse {
    1: required i32 total_clicks,
    2: required list<DailyClickData> daily_clicks,
    3: required list<ReferrerData> referrers,
    4: required list<CountryData> countries,
    5: required list<UserAgentData> user_agents,
}

struct ClickData {
    1: required string ip_address,
    2: required string user_agent,
    3: required string referrer,
    4: required string country,
    5: required string city,
    6: required string clicked_at,
}

struct Pagination {
    1: required i32 page,
    2: required i32 limit,
    3: required i32 total,
    4: required i32 pages,
}

struct ClickHistoryResponse {
    1: required list<ClickData> clicks,
    2: required Pagination pagination,
}

struct GetAnalyticsRequest {
    1: required string code (api.path="code"),
    2: optional string start_date (api.query="start_date"),
    3: optional string end_date (api.query="end_date"),
    4: optional string group_by (api.query="group_by"),
}

struct GetClickHistoryRequest {
    1: required string code (api.path="code"),
    2: optional i32 page (api.query="page"),
    3: optional i32 limit (api.query="limit"),
    4: optional string start_date (api.query="start_date"),
    5: optional string end_date (api.query="end_date"),
}

// Error response
struct ErrorResponse {
    1: required string error,
    2: optional string details,
}

// Service definition
service ShortURLService {
    // Create a new short URL
    CreateShortResponse createShort(1: CreateShortRequest request) (api.post="/api/v1/urls"),
    
    // Get short URL details
    GetShortResponse getShort(1: string code) (api.get="/api/v1/urls/:code"),
    
    // Update short URL
    UpdateShortResponse updateShort(1: UpdateShortRequest request) (api.put="/api/v1/urls/:code"),
    
    // Delete short URL
    DeleteShortResponse deleteShort(1: string code) (api.delete="/api/v1/urls/:code"),
    
    // Get analytics
    AnalyticsResponse getAnalytics(1: GetAnalyticsRequest request) (api.get="/api/v1/urls/:code/analytics"),
    
    // Get click history
    ClickHistoryResponse getClickHistory(1: GetClickHistoryRequest request) (api.get="/api/v1/urls/:code/clicks"),
    
    // Redirect to original URL (returns the original URL)
    string redirect(1: string code) (api.get="/:code"),
}
