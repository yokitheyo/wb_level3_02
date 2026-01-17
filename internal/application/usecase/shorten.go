package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/application/dto"
	"github.com/yokitheyo/wb_level3_02/internal/domain"
)

func (uc *URLShortenerUseCase) Shorten(ctx context.Context, cmd dto.ShortenCommand) (dto.ShortenResult, error) {
	if cmd.URL == "" {
		return dto.ShortenResult{}, ErrURLRequired
	}

	short := cmd.Custom

	if short != "" {
		if err := uc.validator.ValidateCustomShort(short); err != nil {
			return dto.ShortenResult{}, ErrInvalidCustomShort
		}

		existing, err := uc.repo.FindByShort(ctx, short)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("repo error checking custom short")
			return dto.ShortenResult{}, fmt.Errorf("database error")
		}
		if existing != nil {
			return dto.ShortenResult{}, ErrShortCodeAlreadyExists
		}
	} else {
		var err error
		short, err = GenerateShortCode()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("failed to generate short code")
			return dto.ShortenResult{}, fmt.Errorf("failed to generate short code")
		}
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if cmd.Expires > 0 {
		expiresAt = time.Unix(cmd.Expires, 0)
	}

	url := &domain.URL{
		Short:     short,
		Original:  cmd.URL,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Visits:    0,
	}

	if err := uc.repo.Create(ctx, url); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to save URL")
		return dto.ShortenResult{}, fmt.Errorf("failed to save URL")
	}

	if err := uc.cache.Set(ctx, short, cmd.URL, time.Until(expiresAt)); err != nil {
		zlog.Logger.Warn().Err(err).Str("short", short).Msg("failed to cache URL")
	}

	return dto.ShortenResult{
		Short:     short,
		ExpiresAt: expiresAt,
	}, nil
}
