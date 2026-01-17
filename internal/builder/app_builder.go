package builder

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
	wbfretry "github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/api"
	"github.com/yokitheyo/wb_level3_02/internal/application/ports"
	"github.com/yokitheyo/wb_level3_02/internal/application/usecase"
	"github.com/yokitheyo/wb_level3_02/internal/config"
	"github.com/yokitheyo/wb_level3_02/internal/geoip"
	"github.com/yokitheyo/wb_level3_02/internal/infrastructure/geolocation"
	"github.com/yokitheyo/wb_level3_02/internal/infrastructure/repository"
	"github.com/yokitheyo/wb_level3_02/internal/infrastructure/storage"
	"github.com/yokitheyo/wb_level3_02/internal/presentation/handlers"
	internalRetry "github.com/yokitheyo/wb_level3_02/internal/retry"
)

type AppBuilder struct {
	config   *config.Config
	database *dbpg.DB
	redisDB  redis.Client
	retryStr wbfretry.Strategy
	dbOpts   *dbpg.Options

	urlRepo      repository.URLRepository
	urlCache     storage.Cache
	geoIPService ports.GeoService

	urlUseCase *usecase.URLShortenerUseCase
}

func NewAppBuilder(cfg *config.Config) *AppBuilder {
	return &AppBuilder{
		config:   cfg,
		retryStr: internalRetry.DefaultStrategy,
		dbOpts: &dbpg.Options{
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300 * time.Second,
		},
	}
}

func (b *AppBuilder) WithRetryStrategy(strategy wbfretry.Strategy) *AppBuilder {
	if strategy.Attempts > 0 {
		b.retryStr = strategy
	}
	return b
}

func (b *AppBuilder) WithDBOptions(opts *dbpg.Options) *AppBuilder {
	if opts != nil {
		b.dbOpts = opts
	}
	return b
}

func (b *AppBuilder) BuildDatabase() error {
	masterDSN := b.config.Database.DSN
	slaves := []string{}

	var lastErr error
	err := wbfretry.Do(func() error {
		var err error
		b.database, err = dbpg.New(masterDSN, slaves, b.dbOpts)
		if err != nil {
			lastErr = err
			return err
		}
		return nil
	}, b.retryStr)

	if err != nil {
		return fmt.Errorf("failed to connect to database after retries: %w", lastErr)
	}
	return nil
}

func (b *AppBuilder) BuildCache() error {
	rdb := redis.New(
		b.config.Redis.Addr,
		b.config.Redis.Password,
		b.config.Redis.DB,
	)
	b.redisDB = *rdb
	b.urlCache = storage.NewRedisCache(&b.redisDB, "url:", b.retryStr)
	return nil
}

func (b *AppBuilder) BuildMigrations() error {
	if !b.config.Database.Migrate {
		zlog.Logger.Info().Msg("Migrations disabled in config")
		return nil
	}

	migrationsPath := b.config.Database.MigrationsPath
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		zlog.Logger.Warn().Str("path", migrationsPath).Msg("migrations directory not found, skipping")
		return nil
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		b.config.Database.DSN,
	)
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	defer m.Close()

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up error: %w", err)
	}

	zlog.Logger.Info().Msg("Migrations completed successfully")
	return nil
}

func (b *AppBuilder) BuildRepositories() error {
	b.urlRepo = repository.NewPostgresURLRepository(b.database, b.retryStr)
	return nil
}

func (b *AppBuilder) BuildInfrastructure() error {
	geoIPSvc, err := geoip.GetInstance()
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to initialize GeoIP service")
		return nil
	}

	if geoIPSvc != nil {
		b.geoIPService = geolocation.NewGeoIPService(geoIPSvc)
	}
	return nil
}

func (b *AppBuilder) BuildApplicationService() error {
	b.urlUseCase = usecase.NewURLShortenerUseCase(b.urlRepo, b.urlCache, b.geoIPService)
	return nil
}

func (b *AppBuilder) Build() (*api.API, error) {
	if err := b.BuildRepositories(); err != nil {
		return nil, err
	}

	if err := b.BuildInfrastructure(); err != nil {
		return nil, err
	}

	if err := b.BuildApplicationService(); err != nil {
		return nil, err
	}

	handler := handlers.NewURLHandler(b.urlUseCase)

	apiServer := api.NewAPI(handler)

	return apiServer, nil
}
