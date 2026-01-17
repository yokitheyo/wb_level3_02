package usecase

import (
	"context"

	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/application/dto"
	"github.com/yokitheyo/wb_level3_02/internal/domain"
)

func (uc *URLShortenerUseCase) Redirect(ctx context.Context, short string, meta dto.ClickMetadata) (string, error) {
	if short == "" {
		return "", ErrShortCodeRequired
	}

	urlObj, err := uc.repo.FindByShort(ctx, short)
	if err != nil {
		zlog.Logger.Warn().Err(err).Str("short", short).Msg("repo error finding URL")
		return "", ErrNotFound
	}
	if urlObj == nil {
		return "", ErrNotFound
	}

	if err := uc.repo.IncrementVisits(ctx, urlObj.ID); err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to increment visits")
	}

	location := uc.resolveLocation(meta.IP)

	click := &domain.Click{
		URLID:      urlObj.ID,
		Short:      urlObj.Short,
		OccurredAt: urlObj.CreatedAt,
		UserAgent:  meta.UserAgent,
		IP:         location,
		Referrer:   meta.Browser,
		Device:     meta.Device,
	}

	if err := uc.repo.SaveClick(ctx, click); err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to save click")
	}

	return urlObj.Original, nil
}
