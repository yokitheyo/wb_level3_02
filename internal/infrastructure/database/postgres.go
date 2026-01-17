package database

import (
	"github.com/wb-go/wbf/dbpg"
	wbfretry "github.com/wb-go/wbf/retry"
	"github.com/yokitheyo/wb_level3_02/internal/infrastructure/repository"
	"github.com/yokitheyo/wb_level3_02/internal/retry"
)

func MigrateDB(db *dbpg.DB) error {
	// Implementation for running migrations
	return nil
}

func NewPostgresRepository(db *dbpg.DB, strategy wbfretry.Strategy) repository.URLRepository {
	if strategy.Attempts <= 0 {
		strategy = retry.DefaultStrategy
	}
	return repository.NewPostgresURLRepository(db, strategy)
}
