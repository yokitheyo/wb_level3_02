package repository

import (
	"context"
	"time"

	"github.com/yokitheyo/wb_level3_02/internal/domain"
)

type URLRepository interface {
	Create(ctx context.Context, u *domain.URL) error
	FindByShort(ctx context.Context, short string) (*domain.URL, error)
	IncrementVisits(ctx context.Context, id int64) error
	AggregateByDay(ctx context.Context, short string, from, to time.Time) (map[string]int64, error)
	GetDeviceStats(ctx context.Context, short string, from, to time.Time) (map[string]int64, error)
	SaveClick(ctx context.Context, c *domain.Click) error
	GetRecentClicks(ctx context.Context, short string, limit int) ([]*domain.Click, error)
}
