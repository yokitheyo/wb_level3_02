package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	wbfconfig "github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Shortener ShortenerConfig `mapstructure:"shortener"`
}

type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type ShortenerConfig struct {
	BaseURL      string        `mapstructure:"base_url"`
	TTL          time.Duration `mapstructure:"ttl"`
	CleanupEvery time.Duration `mapstructure:"cleanup_every"`
}

func Load(path string) (*Config, error) {
	c := wbfconfig.New()

	if err := c.Load(path); err != nil {
		zlog.Logger.Warn().Err(err).Msg("config file not found, using environment variables and defaults")
	}

	var appConfig Config
	if err := c.Unmarshal(&appConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if appConfig.Server.Addr == "" {
		appConfig.Server.Addr = getEnv("SERVER_ADDR", ":8080")
	}
	if appConfig.Database.DSN == "" {
		appConfig.Database.DSN = getEnv(
			"DB_MASTER",
			"postgres://shortener:shortener@db:5432/shortener?sslmode=disable",
		)
	}
	if appConfig.Redis.Addr == "" {
		appConfig.Redis.Addr = getEnv("REDIS_ADDR", "redis:6379")
	}
	if appConfig.Redis.Password == "" {
		appConfig.Redis.Password = getEnv("REDIS_PASSWORD", "")
	}
	if appConfig.Redis.DB == 0 {
		if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
			if db, err := strconv.Atoi(dbStr); err == nil {
				appConfig.Redis.DB = db
			}
		}
	}
	if appConfig.Shortener.BaseURL == "" {
		appConfig.Shortener.BaseURL = getEnv("BASE_URL", "http://localhost:8080")
	}
	if appConfig.Shortener.TTL == 0 {
		appConfig.Shortener.TTL = 24 * time.Hour
	}
	if appConfig.Shortener.CleanupEvery == 0 {
		appConfig.Shortener.CleanupEvery = time.Hour
	}

	zlog.Logger.Info().Msg("configuration loaded")
	return &appConfig, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
