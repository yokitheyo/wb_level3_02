package dto

import "time"

// ShortenCommand - request to shorten URL
type ShortenCommand struct {
	URL     string
	Custom  string
	Expires int64 // unix timestamp
}

// ShortenResult - result of URL shortening
type ShortenResult struct {
	Short     string
	ExpiresAt time.Time
}

// RedirectCommand - request to resolve and redirect
type RedirectCommand struct {
	Short     string
	IP        string
	UserAgent string
}

// ClickMetadata - metadata for recording a click
type ClickMetadata struct {
	IP        string
	UserAgent string
	Location  string
	Browser   string
	Device    string
	Referrer  string
}

// AnalyticsQuery - query for URL analytics
type AnalyticsQuery struct {
	Short string
}

// AnalyticsResult - analytics for a URL
type AnalyticsResult struct {
	Short      string
	Original   string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	VisitCount int64
}

// DetailedAnalyticsQuery - query for detailed analytics
type DetailedAnalyticsQuery struct {
	Short string
	From  time.Time
	To    time.Time
}

// DetailedAnalyticsResult - detailed analytics breakdown
type DetailedAnalyticsResult struct {
	Short            string
	DailyClicks      map[string]int64
	DeviceStats      map[string]int64
	MobilePercentage float64
	TotalClicks      int64
}

// RecentClicksQuery - query for recent clicks
type RecentClicksQuery struct {
	Short string
	Limit int
}

// RecentClicksResult - recent clicks data
type RecentClicksResult struct {
	Clicks []map[string]interface{}
	Total  int
}
