package dto

// ShortenRequest - HTTP POST /shorten request body
type ShortenRequest struct {
	URL     string `json:"url" binding:"required"`
	Custom  string `json:"custom"`
	Expires int64  `json:"expires"`
}

// AnalyticsDetailedQueryParams - HTTP query parameters for detailed analytics
type AnalyticsDetailedQueryParams struct {
	From string `form:"from"` // Format: 2006-01-02
	To   string `form:"to"`   // Format: 2006-01-02
}

// RecentClicksQueryParams - HTTP query parameters for recent clicks
type RecentClicksQueryParams struct {
	Limit string `form:"limit"`
}
