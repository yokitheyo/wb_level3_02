package geolocation

type Service interface {
	GetLocationFromIP(ip string) string
}
