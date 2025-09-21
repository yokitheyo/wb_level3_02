package model

import "time"

type Click struct {
	ID        int64     `json:"id" db:"id"`
	URLID     int64     `json:"url_id" db:"url_id"`
	Short     string    `json:"short" db:"short"`
	Occurred  time.Time `json:"occurred_at" db:"occurred_at"`
	UserAgent string    `json:"user_agent,omitempty" db:"user_agent"`
	IP        string    `json:"ip,omitempty" db:"ip"`
	Referrer  string    `json:"referrer,omitempty" db:"referrer"`
	Device    string    `json:"device,omitempty" db:"device"`
}
