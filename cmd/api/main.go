package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"

	"github.com/yokitheyo/wb_level3_02/internal/api"
	"github.com/yokitheyo/wb_level3_02/internal/app"
	"github.com/yokitheyo/wb_level3_02/internal/cache"
	"github.com/yokitheyo/wb_level3_02/internal/config"
	"github.com/yokitheyo/wb_level3_02/internal/db"
	"github.com/yokitheyo/wb_level3_02/internal/repo"
	internalRetry "github.com/yokitheyo/wb_level3_02/internal/retry"
	"github.com/yokitheyo/wb_level3_02/internal/service"
)

func main() {
	zlog.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	configPath := "config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "/app/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load config")
	}

	masterDSN := cfg.Database.DSN
	slaves := []string{}

	dbOpts := &dbpg.Options{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 300 * time.Second,
	}

	var database *dbpg.DB
	for i := 0; i < 10; i++ {
		database, err = dbpg.New(masterDSN, slaves, dbOpts)
		if err == nil {
			break
		}
		zlog.Logger.Warn().Err(err).Msg("waiting for database to be ready")
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to database after retries")
	}

	if err := db.RunMigrations(database, "/app/migrations"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("migrations failed")
	}

	repo := repo.NewPostgresRepo(database, internalRetry.DefaultStrategy)

	rdb := redis.New(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	cache := cache.NewRedisCache(rdb, "url:", internalRetry.DefaultStrategy)

	svc := service.NewURLService(repo, cache)
	app := app.NewApp(repo, cache, svc)

	apiServer := api.NewAPI(app)
	go func() {
		if err := apiServer.Start(cfg.Server.Addr); err != nil && err != http.ErrServerClosed {
			zlog.Logger.Fatal().Err(err).Msg("failed to start API server")
		}
	}()

	<-ctx.Done()
	apiServer.Stop(ctx)
}
