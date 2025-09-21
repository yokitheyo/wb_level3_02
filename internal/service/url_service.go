package service

import (
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/cache"
	"github.com/yokitheyo/wb_level3_02/internal/model"
	"github.com/yokitheyo/wb_level3_02/internal/repo"
	"github.com/yokitheyo/wb_level3_02/internal/util"
)

type URLService struct {
	repo  *repo.PostgresRepo
	cache *cache.RedisCache
}

func NewURLService(repo *repo.PostgresRepo, cache *cache.RedisCache) *URLService {
	return &URLService{repo, cache}
}

func GenerateShortCode() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(uuid.String()))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(encoded, "=")[:6], nil
}

func (s *URLService) HandleShorten(c *ginext.Context) {
	type request struct {
		URL     string `json:"url" binding:"required"`
		Expires int64  `json:"expires"`
	}

	var req request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	short, err := GenerateShortCode()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to generate short code")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "internal error"})
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if req.Expires > 0 {
		expiresAt = time.Unix(req.Expires, 0)
	}

	url := &model.URL{
		Short:     short,
		Original:  req.URL,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Visits:    0,
	}

	if err := s.repo.Create(ctx, url); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to save URL to DB")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "internal error"})
		return
	}

	if err := s.cache.Set(ctx, short, req.URL, time.Until(expiresAt)); err != nil {
		zlog.Logger.Warn().Err(err).Str("short", short).Msg("failed to save URL in Redis")
	}

	c.JSON(http.StatusOK, ginext.H{
		"short":   short,
		"expires": expiresAt.Unix(),
	})
}

func (s *URLService) HandleRedirect(c *ginext.Context) {
	short := c.Param("short")
	ctx := c.Request.Context()

	var original string
	var urlObj *model.URL

	original, err := s.cache.Get(ctx, short)
	if err != nil || original == "" {
		u, err := s.repo.FindByShort(ctx, short)
		if err != nil {
			zlog.Logger.Warn().Err(err).Str("short", short).Msg("url not found")
			c.JSON(http.StatusNotFound, ginext.H{"error": "url not found"})
			return
		}
		urlObj = u
		original = u.Original

		_ = s.cache.Set(ctx, short, original, time.Until(u.ExpiresAt))
	}

	if urlObj == nil {
		if u, err := s.repo.FindByShort(ctx, short); err == nil {
			urlObj = u
		}
	}

	var click *model.Click
	if urlObj != nil {
		if err := s.repo.IncrementVisits(ctx, urlObj.ID); err != nil {
			zlog.Logger.Warn().Int64("url_id", urlObj.ID).Err(err).Msg("failed to increment visits")
		}

		click = &model.Click{
			URLID:     urlObj.ID,
			Short:     urlObj.Short,
			Occurred:  time.Now().UTC(),
			UserAgent: c.Request.UserAgent(),
			IP:        c.ClientIP(),
			Referrer:  c.Request.Referer(),
			Device:    util.DetectDevice(c.Request.UserAgent()),
		}
	}

	if click != nil {
		if err := s.repo.SaveClick(ctx, click); err != nil {
			zlog.Logger.Warn().Err(err).Msg("failed to save click")
		}
	}

	c.Redirect(http.StatusFound, original)
}

func (s *URLService) HandleAnalytics(c *ginext.Context) {
	short := c.Param("short")
	ctx := c.Request.Context()

	u, err := s.repo.FindByShort(ctx, short)
	if err != nil {
		zlog.Logger.Warn().Err(err).Str("short", short).Msg("URL not found for analytics")
		c.JSON(http.StatusNotFound, ginext.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"short":       u.Short,
		"original":    u.Original,
		"created_at":  u.CreatedAt.Unix(),
		"expires_at":  u.ExpiresAt.Unix(),
		"visit_count": u.Visits,
	})
}

func (s *URLService) HandleDetailedAnalytics(c *ginext.Context) {
	short := c.Param("short")
	ctx := c.Request.Context()

	fromStr := c.Query("from")
	toStr := c.Query("to")

	to := time.Now()
	from := to.AddDate(0, 0, -30)

	if fromStr != "" {
		if parsedFrom, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = parsedFrom
		}
	}

	if toStr != "" {
		if parsedTo, err := time.Parse("2006-01-02", toStr); err == nil {
			to = parsedTo.Add(24 * time.Hour)
		}
	}

	u, err := s.repo.FindByShort(ctx, short)
	if err != nil || u == nil {
		zlog.Logger.Warn().Err(err).Str("short", short).Msg("URL not found for detailed analytics")
		c.JSON(http.StatusNotFound, ginext.H{"error": "not found"})
		return
	}

	dailyClicks, err := s.repo.AggregateByDay(ctx, short, from, to)
	if err != nil {
		zlog.Logger.Warn().Err(err).Str("short", short).Msg("failed to get daily clicks")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "internal error"})
		return
	}

	deviceStats, err := s.repo.GetDeviceStats(ctx, short, from, to)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to get device stats")
		deviceStats = make(map[string]int64)
	}

	mobilePercent := 0.0
	totalClicks := u.Visits
	if totalClicks > 0 && deviceStats["mobile"] > 0 {
		mobilePercent = float64(deviceStats["mobile"]) / float64(totalClicks) * 100
	}

	c.JSON(http.StatusOK, ginext.H{
		"short":             u.Short,
		"daily_clicks":      dailyClicks,
		"device_stats":      deviceStats,
		"mobile_percentage": int(mobilePercent),
		"total_clicks":      totalClicks,
	})
}

func (s *URLService) HandleRecentClicks(c *ginext.Context) {
	short := c.Param("short")
	ctx := c.Request.Context()

	u, err := s.repo.FindByShort(ctx, short)
	if err != nil || u == nil {
		zlog.Logger.Warn().Err(err).Str("short", short).Msg("URL not found for recent clicks")
		c.JSON(http.StatusNotFound, ginext.H{"error": "not found"})
		return
	}

	clicks, err := s.repo.GetRecentClicks(ctx, short, 50)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get recent clicks")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"clicks": clicks,
		"total":  len(clicks),
	})
}
