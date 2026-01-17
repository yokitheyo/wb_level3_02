package middleware

import (
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// PanicRecoveryMiddleware recovers from panics
func PanicRecoveryMiddleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		defer func() {
			if err := recover(); err != nil {
				zlog.Logger.Error().
					Interface("panic", err).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("Panic recovered in middleware")
				c.JSON(500, ginext.H{"error": "internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RequestLoggingMiddleware logs request details
func RequestLoggingMiddleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		zlog.Logger.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", duration).
			Msg("HTTP request")
	}
}

// CORSMiddleware adds CORS headers
func CORSMiddleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
