package app

import (
	"github.com/yokitheyo/wb_level3_02/internal/cache"
	"github.com/yokitheyo/wb_level3_02/internal/repo"
	"github.com/yokitheyo/wb_level3_02/internal/service"
)

type App struct {
	Repo    *repo.PostgresRepo
	Cache   *cache.RedisCache
	Service *service.URLService
}

func NewApp(repo *repo.PostgresRepo, cache *cache.RedisCache, service *service.URLService) *App {
	return &App{
		Repo:    repo,
		Cache:   cache,
		Service: service,
	}
}
