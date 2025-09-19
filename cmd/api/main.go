package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/app"
	"github.com/yokitheyo/wb_level3_02/internal/cache"
	"github.com/yokitheyo/wb_level3_02/internal/repo"
	internalRetry "github.com/yokitheyo/wb_level3_02/internal/retry"
	"github.com/yokitheyo/wb_level3_02/internal/service"
)

func main() {
	zlog.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.New()
	if err := cfg.Load("config.yaml"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load config")
	}

	masterDSN := cfg.GetString("database.master")
	dbOpts := &dbpg.Options{
		MaxOpenConns:    cfg.GetInt("database.max_open_conns"),
		MaxIdleConns:    cfg.GetInt("database.max_idle_conns"),
		ConnMaxLifetime: time.Duration(cfg.GetInt("database.conn_max_lifetime_sec")) * time.Second,
	}
	db, err := dbpg.New(masterDSN, nil, dbOpts)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	rdb := redis.New(
		cfg.GetString("redis.addr"),
		cfg.GetString("redis.password"),
		cfg.GetInt("redis.db"),
	)
	cache := cache.NewRedisCache(rdb, "url:", internalRetry.DefaultStrategy)

	repo := repo.NewPostgresRepo(db, internalRetry.DefaultStrategy)
	svc := service.NewURLService(repo, cache)

	a := app.NewApp(repo, cache, svc)

	go func() {
		if err := a.Start(cfg.GetString("server.addr")); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to start app")
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Stop(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to stop server gracefully")
	} else {
		zlog.Logger.Info().Msg("server stopped gracefully")
	}
}
