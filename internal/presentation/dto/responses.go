package dto

type ShortenResponse struct {
	Short   string `json:"short"`
	Expires int64  `json:"expires"`
}

type AnalyticsResponse struct {
	Short      string `json:"short"`
	Original   string `json:"original"`
	CreatedAt  int64  `json:"created_at"`
	ExpiresAt  int64  `json:"expires_at"`
	VisitCount int64  `json:"visit_count"`
}

type DetailedAnalyticsResponse struct {
	Short            string           `json:"short"`
	DailyClicks      map[string]int64 `json:"daily_clicks"`
	DeviceStats      map[string]int64 `json:"device_stats"`
	MobilePercentage int              `json:"mobile_percentage"`
	TotalClicks      int64            `json:"total_clicks"`
}

type RecentClicksResponse struct {
	Clicks []ClickData `json:"clicks"`
	Total  int         `json:"total"`
}

type ClickData struct {
	OccurredAt string `json:"occurred_at"`
	IP         string `json:"ip"`
	Referrer   string `json:"referrer"`
	Device     string `json:"device"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
