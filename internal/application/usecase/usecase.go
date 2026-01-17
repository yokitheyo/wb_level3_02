package usecase

import (
	"github.com/yokitheyo/wb_level3_02/internal/application/ports"
)

type URLShortenerUseCase struct {
	repo       ports.URLRepository
	cache      ports.Cache
	geoService ports.GeoService
	validator  *ShortCodeValidator
}

func NewURLShortenerUseCase(
	repo ports.URLRepository,
	cache ports.Cache,
	geoService ports.GeoService,
) *URLShortenerUseCase {
	return &URLShortenerUseCase{
		repo:       repo,
		cache:      cache,
		geoService: geoService,
		validator:  NewShortCodeValidator(),
	}
}
