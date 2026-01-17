package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/wb-go/wbf/zlog"

	"github.com/yokitheyo/wb_level3_02/internal/builder"
	"github.com/yokitheyo/wb_level3_02/internal/config"
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

	appBuilder := builder.NewAppBuilder(cfg)

	if err := appBuilder.BuildDatabase(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to database after retries")
	}

	if err := appBuilder.BuildMigrations(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("migrations failed")
	}

	if err := appBuilder.BuildCache(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to initialize cache")
	}

	apiServer, err := appBuilder.Build()
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to build application")
	}

	go func() {
		if err := apiServer.Start(cfg.Server.Addr); err != nil && err != http.ErrServerClosed {
			zlog.Logger.Fatal().Err(err).Msg("failed to start API server")
		}
	}()

	<-ctx.Done()
	apiServer.Stop(ctx)
}
