package model

import "time"

type URL struct {
	ID        int64     `json:"id"`
	Short     string    `json:"short"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Visits    int64     `json:"visits"`
}
