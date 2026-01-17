package api

import (
	"context"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"github.com/yokitheyo/wb_level3_02/internal/presentation/handlers"
	"github.com/yokitheyo/wb_level3_02/internal/presentation/router"
)

type API struct {
	Engine  *ginext.Engine
	server  *http.Server
	handler *handlers.URLHandler
	router  *router.Router
}

func NewAPI(handler *handlers.URLHandler) *API {
	engine := ginext.New("")

	api := &API{
		Engine:  engine,
		handler: handler,
	}

	api.router = router.NewRouter(engine, handler)
	api.router.Setup()

	return api
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
