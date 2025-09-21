package api

import (
	"context"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/app"
)

type API struct {
	Engine *ginext.Engine
	server *http.Server
	App    *app.App
}

func NewAPI(a *app.App) *API {
	engine := ginext.New()
	engine.Use(ginext.Logger(), ginext.Recovery())
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/static", "./static")
	engine.GET("/", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := &API{
		Engine: engine,
		App:    a,
	}

	api.registerRoutes()
	return api
}

func (api *API) registerRoutes() {
	api.Engine.POST("/shorten", api.App.Service.HandleShorten)
	api.Engine.GET("/s/:short", api.App.Service.HandleRedirect)
	api.Engine.GET("/analytics/:short", api.App.Service.HandleAnalytics)
}

func (api *API) Start(addr string) error {
	api.server = &http.Server{
		Addr:    addr,
		Handler: api.Engine,
	}

	zlog.Logger.Info().Msgf("Starting server on %s", addr)
	return api.server.ListenAndServe()
}

func (api *API) Stop(ctx context.Context) error {
	if api.server == nil {
		return nil
	}

	zlog.Logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return api.server.Shutdown(ctx)
}
