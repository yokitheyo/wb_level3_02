package repository

import (
	"github.com/yokitheyo/wb_level3_02/internal/application/ports"
)

// PostgresURLRepository implements ports.URLRepository interface
var _ ports.URLRepository = (*PostgresURLRepository)(nil)
