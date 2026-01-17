package handlers

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
	"github.com/yokitheyo/wb_level3_02/internal/application/dto"
	"github.com/yokitheyo/wb_level3_02/internal/application/usecase"
	"github.com/yokitheyo/wb_level3_02/internal/presentation"
	presentationdto "github.com/yokitheyo/wb_level3_02/internal/presentation/dto"
	presentationutil "github.com/yokitheyo/wb_level3_02/internal/presentation/util"
)

type URLHandler struct {
	useCase *usecase.URLShortenerUseCase
}

func NewURLHandler(uc *usecase.URLShortenerUseCase) *URLHandler {
	return &URLHandler{
		useCase: uc,
	}
}

// HandleShorten - HTTP POST /shorten
func (h *URLHandler) HandleShorten(c *ginext.Context) {
	var req presentationdto.ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	cmd := dto.ShortenCommand{
		URL:     req.URL,
		Custom:  req.Custom,
		Expires: req.Expires,
	}

	result, err := h.useCase.Shorten(c.Request.Context(), cmd)
	if err != nil {
		status := presentationutil.MapErrorToStatus(err)
		c.JSON(status, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"short":   result.Short,
		"expires": result.ExpiresAt.Unix(),
	})
}

// HandleRedirect - HTTP GET /s/:short
func (h *URLHandler) HandleRedirect(c *ginext.Context) {
	short := c.Param("short")

	meta := presentationutil.BuildClickMetadata(c)

	originalURL, err := h.useCase.Redirect(c.Request.Context(), short, meta)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

// HandleAnalytics - HTTP GET /analytics/:short
func (h *URLHandler) HandleAnalytics(c *ginext.Context) {
	short := c.Param("short")

	query := dto.AnalyticsQuery{Short: short}

	result, err := h.useCase.GetAnalytics(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"short":       result.Short,
		"original":    result.Original,
		"created_at":  result.CreatedAt.Unix(),
		"expires_at":  result.ExpiresAt.Unix(),
		"visit_count": result.VisitCount,
	})
}

// HandleDetailedAnalytics - HTTP GET /analytics/:short/detailed
func (h *URLHandler) HandleDetailedAnalytics(c *ginext.Context) {
	short := c.Param("short")

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, to := presentationutil.ParseDateRange(fromStr, toStr)

	query := dto.DetailedAnalyticsQuery{
		Short: short,
		From:  from,
		To:    to,
	}

	result, err := h.useCase.GetDetailedAnalytics(c.Request.Context(), query)
	if err != nil {
		status := presentationutil.MapErrorToStatus(err)
		c.JSON(status, ginext.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"short":             result.Short,
		"daily_clicks":      result.DailyClicks,
		"device_stats":      result.DeviceStats,
		"mobile_percentage": int(result.MobilePercentage),
		"total_clicks":      result.TotalClicks,
	})
}

// HandleRecentClicks - HTTP GET /analytics/:short/recent-clicks
func (h *URLHandler) HandleRecentClicks(c *ginext.Context) {
	short := c.Param("short")

	limitStr := c.Query("limit")
	limit := presentationutil.ParseLimit(limitStr, presentation.DefaultRecentClicksLimit, presentation.MaxRecentClicksLimit)

	query := dto.RecentClicksQuery{
		Short: short,
		Limit: limit,
	}

	result, err := h.useCase.GetRecentClicks(c.Request.Context(), query)
	if err != nil {
		status := presentationutil.MapErrorToStatus(err)
		c.JSON(status, ginext.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"clicks": result.Clicks,
		"total":  result.Total,
	})
}
