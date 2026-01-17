package geolocation

import (
	"strings"

	"github.com/yokitheyo/wb_level3_02/internal/geoip"
	"github.com/yokitheyo/wb_level3_02/internal/util"
)

type GeoIPServiceImpl struct {
	svc *geoip.Service
}

func NewGeoIPService(svc *geoip.Service) Service {
	return &GeoIPServiceImpl{svc: svc}
}

func (a *GeoIPServiceImpl) GetLocationFromIP(ip string) string {
	if strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "172.") || strings.HasPrefix(ip, "10.") {
		return "Local Network"
	}

	if a.svc != nil {
		return a.svc.GetLocationFromIP(ip)
	}
	return "Unknown"
}

func GetBrowserFromUserAgent(ua string) string {
	return geoip.GetBrowserFromUserAgent(ua)
}

func DetectDevice(ua string) string {
	return util.DetectDevice(ua)
}
