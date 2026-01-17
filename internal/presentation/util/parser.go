package util

import (
	"strconv"
	"time"

	"github.com/yokitheyo/wb_level3_02/internal/presentation"
)

func ParseDateRange(fromStr, toStr string) (time.Time, time.Time) {
	to := time.Now()
	from := to.AddDate(0, 0, -presentation.DefaultAnalyticsDaysBack)

	if fromStr != "" {
		if parsed, err := time.Parse(presentation.DateFormat, fromStr); err == nil {
			from = parsed
		}
	}

	if toStr != "" {
		if parsed, err := time.Parse(presentation.DateFormat, toStr); err == nil {
			to = parsed.Add(24 * time.Hour)
		}
	}

	return from, to
}

func ParseLimit(limitStr string, defaultLimit, maxLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return defaultLimit
	}

	if limit > maxLimit {
		return maxLimit
	}

	return limit
}
