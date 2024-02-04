package models

import "time"

type Session struct {
	RefreshToken string
	ExpiresAt    time.Time
}
