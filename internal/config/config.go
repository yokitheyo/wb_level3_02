package config

import (
	"fmt"
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
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var appConfig Config
	if err := c.Unmarshal(&appConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if appConfig.Server.Addr == "" {
		appConfig.Server.Addr = ":8080"
	}
	if appConfig.Shortener.TTL == 0 {
		appConfig.Shortener.TTL = 24 * time.Hour
	}
	if appConfig.Shortener.CleanupEvery == 0 {
		appConfig.Shortener.CleanupEvery = time.Hour
	}

	zlog.Logger.Info().Str("config_path", path).Msg("configuration loaded")
	return &appConfig, nil
}
