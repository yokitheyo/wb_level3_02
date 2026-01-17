package usecase

import (
	"context"
	"fmt"

	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/application/dto"
)

func (uc *URLShortenerUseCase) GetAnalytics(ctx context.Context, query dto.AnalyticsQuery) (dto.AnalyticsResult, error) {
	if query.Short == "" || len(query.Short) > 50 {
		return dto.AnalyticsResult{}, ErrInvalidQuery
	}

	u, err := uc.repo.FindByShort(ctx, query.Short)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("repo error in analytics")
		return dto.AnalyticsResult{}, ErrNotFound
	}
	if u == nil {
		return dto.AnalyticsResult{}, ErrNotFound
	}

	return dto.AnalyticsResult{
		Short:      u.Short,
		Original:   u.Original,
		CreatedAt:  u.CreatedAt,
		ExpiresAt:  u.ExpiresAt,
		VisitCount: u.Visits,
	}, nil
}

func (uc *URLShortenerUseCase) GetDetailedAnalytics(ctx context.Context, query dto.DetailedAnalyticsQuery) (dto.DetailedAnalyticsResult, error) {
	if query.Short == "" {
		return dto.DetailedAnalyticsResult{}, fmt.Errorf("short code is required")
	}

	u, err := uc.repo.FindByShort(ctx, query.Short)
	if err != nil || u == nil {
		return dto.DetailedAnalyticsResult{}, fmt.Errorf("not found")
	}

	dailyClicks, err := uc.repo.AggregateByDay(ctx, query.Short, query.From, query.To)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to get daily stats")
		dailyClicks = make(map[string]int64)
	}

	deviceStats, err := uc.repo.GetDeviceStats(ctx, query.Short, query.From, query.To)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to get device stats")
		deviceStats = make(map[string]int64)
	}

	mobilePercentage := 0.0
	if u.Visits > 0 && deviceStats["mobile"] > 0 {
		mobilePercentage = float64(deviceStats["mobile"]) / float64(u.Visits) * 100
	}

	return dto.DetailedAnalyticsResult{
		Short:            u.Short,
		DailyClicks:      dailyClicks,
		DeviceStats:      deviceStats,
		MobilePercentage: mobilePercentage,
		TotalClicks:      u.Visits,
	}, nil
}
