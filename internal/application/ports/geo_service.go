package ports

// GeoService resolves geographic location from IP
type GeoService interface {
	GetLocationFromIP(ip string) string
}
