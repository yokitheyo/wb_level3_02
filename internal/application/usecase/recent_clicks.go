package usecase

import (
	"context"
	"fmt"

	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/application/dto"
	"github.com/yokitheyo/wb_level3_02/internal/domain"
)

func (uc *URLShortenerUseCase) GetRecentClicks(ctx context.Context, query dto.RecentClicksQuery) (dto.RecentClicksResult, error) {
	if query.Short == "" {
		return dto.RecentClicksResult{}, ErrShortCodeRequired
	}

	u, err := uc.repo.FindByShort(ctx, query.Short)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("repo error in recent clicks")
		return dto.RecentClicksResult{}, ErrNotFound
	}
	if u == nil {
		return dto.RecentClicksResult{}, ErrNotFound
	}

	limit := query.Limit
	if limit <= 0 || limit > 1000 {
		limit = 50
	}

	clicks, err := uc.repo.GetRecentClicks(ctx, query.Short, limit)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get recent clicks")
		return dto.RecentClicksResult{}, fmt.Errorf("failed to get recent clicks")
	}

	if clicks == nil {
		clicks = []*domain.Click{}
	}

	clickMaps := make([]map[string]interface{}, len(clicks))
	for i, click := range clicks {
		clickMaps[i] = map[string]interface{}{
			"url_id":      click.URLID,
			"short":       click.Short,
			"occurred_at": click.OccurredAt.UnixMilli(),
			"user_agent":  click.UserAgent,
			"ip":          click.IP,
			"referrer":    click.Referrer,
			"device":      click.Device,
		}
	}

	return dto.RecentClicksResult{
		Clicks: clickMaps,
		Total:  len(clickMaps),
	}, nil
}
