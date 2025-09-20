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

	cfg, err := config.Load("config.yaml")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load config")
	}

	masterDSN := cfg.Database.DSN

	slaves := []string{}

	dbOpts := &dbpg.Options{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Duration(300) * time.Second,
	}

	database, err := dbpg.New(masterDSN, slaves, dbOpts)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	if err := db.RunMigrations(database, "file://migrations"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to run migrations")
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
