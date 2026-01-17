package geolocation

import (
	"github.com/yokitheyo/wb_level3_02/internal/application/ports"
)

// GeoIPServiceImpl implements ports.GeoService interface
var _ ports.GeoService = (*GeoIPServiceImpl)(nil)
