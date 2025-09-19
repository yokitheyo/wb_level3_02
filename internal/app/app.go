package app

import (
	"context"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/cache"
	"github.com/yokitheyo/wb_level3_02/internal/repo"
	"github.com/yokitheyo/wb_level3_02/internal/service"
)

type App struct {
	Repo    *repo.PostgresRepo
	Cache   *cache.RedisCache
	Service *service.URLService
	Engine  *ginext.Engine
	server  *http.Server
}

func NewApp(repo *repo.PostgresRepo, cache *cache.RedisCache, service *service.URLService) *App {
	engine := ginext.New()

	engine.Use(ginext.Logger(), ginext.Recovery())
	engine.LoadHTMLGlob("templates/*")

	app := &App{
		Repo:    repo,
		Cache:   cache,
		Service: service,
		Engine:  engine,
	}

	app.registerRoutes()
	return app
}

func (a *App) registerRoutes() {
	a.Engine.POST("/shorten", a.Service.HandleShorten)
	a.Engine.GET("/s/:short", a.Service.HandleRedirect)
	a.Engine.GET("/analytics/:short", a.Service.HandleAnalytics)
}

func (a *App) Start(addr string) error {
	a.server = &http.Server{
		Addr:    addr,
		Handler: a.Engine,
	}

	zlog.Logger.Info().Msgf("Starting server on %s", addr)
	return a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	if a.server == nil {
		return nil
	}

	zlog.Logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.server.Shutdown(ctx)
}
