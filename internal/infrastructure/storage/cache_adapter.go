package storage

import (
	"github.com/yokitheyo/wb_level3_02/internal/application/ports"
)

// Ensure RedisCache implements ports.Cache interface
var _ ports.Cache = (*RedisCache)(nil)
