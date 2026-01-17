package domain

import (
	"time"
)

type URL struct {
	ID        int64
	Short     string
	Original  string
	CreatedAt time.Time
	ExpiresAt time.Time
	Visits    int64
}

type Click struct {
	ID         int64
	URLID      int64
	Short      string
	OccurredAt time.Time
	UserAgent  string
	IP         string
	Referrer   string
	Device     string
}
