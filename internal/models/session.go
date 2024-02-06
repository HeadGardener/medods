package models

import "time"

type Session struct {
	UserID       string    `bson:"user_id"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiresAt    time.Time `bson:"expires_at"`
}
