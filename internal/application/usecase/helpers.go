package usecase

func (uc *URLShortenerUseCase) resolveLocation(ip string) string {
	if isLocalIP(ip) {
		return "Moscow, Moscow, Russia"
	}

	if uc.geoService != nil {
		loc := uc.geoService.GetLocationFromIP(ip)
		if loc != "" {
			return loc
		}
	}

	return "Unknown"
}

func isLocalIP(ip string) bool {
	if ip == "" {
		return true
	}
	return ip[0:4] == "127." || ip[0:4] == "10." || ip[0:6] == "172.1"
}
