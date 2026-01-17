package router

import (
	"github.com/wb-go/wbf/ginext"
	"github.com/yokitheyo/wb_level3_02/internal/presentation/handlers"
	"github.com/yokitheyo/wb_level3_02/internal/presentation/middleware"
)

type Router struct {
	engine  *ginext.Engine
	handler *handlers.URLHandler
}

func NewRouter(engine *ginext.Engine, handler *handlers.URLHandler) *Router {
	return &Router{
		engine:  engine,
		handler: handler,
	}
}

func (r *Router) Setup() {
	r.engine.Use(middleware.PanicRecoveryMiddleware())
	r.engine.Use(middleware.RequestLoggingMiddleware())
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(ginext.Logger(), ginext.Recovery())

	r.engine.LoadHTMLGlob("templates/*")
	r.engine.Static("/static", "./static")

	r.engine.GET("/", func(c *ginext.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.engine.POST("/shorten", r.handler.HandleShorten)
	r.engine.GET("/s/:short", r.handler.HandleRedirect)

	r.engine.GET("/analytics/:short", r.handler.HandleAnalytics)
	r.engine.GET("/analytics/:short/detailed", r.handler.HandleDetailedAnalytics)
	r.engine.GET("/analytics/:short/recent-clicks", r.handler.HandleRecentClicks)
}
