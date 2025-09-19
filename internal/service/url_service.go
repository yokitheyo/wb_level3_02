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

	if urlObj != nil {
		if err := s.repo.IncrementVisits(ctx, urlObj.ID); err != nil {
			zlog.Logger.Warn().Int64("url_id", urlObj.ID).Err(err).Msg("failed to increment visits")
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
