package model

import "time"

type Click struct {
	ID        int64     `json:"id"`
	URLID     int64     `json:"url_id"`
	Short     string    `json:"short"`
	Occurred  time.Time `json:"occurred_at"`
	UserAgent string    `json:"user_agent,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Referrer  string    `json:"referrer,omitempty"`
	Device    string    `json:"device,omitempty"`
}
