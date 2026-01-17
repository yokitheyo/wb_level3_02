package config

import (
	"fmt"
	"os"
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
	DSN            string `mapstructure:"dsn"`
	Migrate        bool   `mapstructure:"migrate"`
	MigrationsPath string `mapstructure:"migrations_path"`
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

	if err := c.LoadConfigFiles(path); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var appConfig Config
	if err := c.Unmarshal(&appConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with environment variables if set
	if dsn := os.Getenv("DB_MASTER"); dsn != "" {
		appConfig.Database.DSN = dsn
	}
	if addr := os.Getenv("SERVER_ADDR"); addr != "" {
		appConfig.Server.Addr = addr
	}
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		appConfig.Redis.Addr = addr
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		appConfig.Redis.Password = password
	}

	zlog.Logger.Info().Msg("configuration loaded")
	return &appConfig, nil
}
