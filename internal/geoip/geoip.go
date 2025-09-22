// TODO: Fix location detection for local IP addresses (192.168.x.x, etc.).
package geoip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Service struct {
	apiKey string
	client *http.Client
}

type AbstractGeoIPResponse struct {
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

var (
	instance *Service
	once     sync.Once
)

// GetInstance инициализация singleton
func GetInstance() (*Service, error) {
	var err error
	once.Do(func() {
		instance = &Service{
			apiKey: "07eed3774e8545ad91d3d21ae0ae204b",
			client: &http.Client{},
		}
	})
	return instance, err
}

func (s *Service) GetLocationFromIP(ip string) string {
	if ip == "" {
		return "Unknown"
	}

	// На случай локальных IP всё равно идем на AbstractAPI
	url := fmt.Sprintf("https://ipgeolocation.abstractapi.com/v1/?api_key=%s&ip_address=%s", s.apiKey, ip)
	resp, err := s.client.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Unknown"
	}

	var geo AbstractGeoIPResponse
	if err := json.Unmarshal(body, &geo); err != nil {
		return "Unknown"
	}

	parts := []string{}
	if geo.City != "" {
		parts = append(parts, geo.City)
	}
	if geo.Region != "" {
		parts = append(parts, geo.Region)
	}
	if geo.Country != "" {
		parts = append(parts, geo.Country)
	}

	if len(parts) == 0 {
		return "Unknown"
	}

	return strings.Join(parts, ", ")
}

func GetBrowserFromUserAgent(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	if strings.Contains(userAgent, "edg/") {
		return "Microsoft Edge"
	} else if strings.Contains(userAgent, "chrome/") && !strings.Contains(userAgent, "chromium") && !strings.Contains(userAgent, "edg") {
		if strings.Contains(userAgent, "yabrowser") {
			return "Yandex Browser"
		}
		return "Google Chrome"
	} else if strings.Contains(userAgent, "firefox/") {
		return "Mozilla Firefox"
	} else if strings.Contains(userAgent, "safari/") && !strings.Contains(userAgent, "chrome") {
		return "Safari"
	} else if strings.Contains(userAgent, "opera") || strings.Contains(userAgent, "opr/") {
		return "Opera"
	} else if strings.Contains(userAgent, "msie") || strings.Contains(userAgent, "trident") {
		return "Internet Explorer"
	} else if strings.Contains(userAgent, "chromium") {
		return "Chromium"
	}

	return "Unknown Browser"
}
