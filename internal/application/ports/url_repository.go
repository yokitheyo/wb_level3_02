package ports

import (
	"context"
	"time"

	"github.com/yokitheyo/wb_level3_02/internal/domain"
)

// URLRepository defines domain.URL persistence operations
type URLRepository interface {
	Create(ctx context.Context, url *domain.URL) error
	FindByShort(ctx context.Context, short string) (*domain.URL, error)
	IncrementVisits(ctx context.Context, id int64) error
	SaveClick(ctx context.Context, click *domain.Click) error
	AggregateByDay(ctx context.Context, short string, from, to time.Time) (map[string]int64, error)
	GetDeviceStats(ctx context.Context, short string, from, to time.Time) (map[string]int64, error)
	GetRecentClicks(ctx context.Context, short string, limit int) ([]*domain.Click, error)
}
