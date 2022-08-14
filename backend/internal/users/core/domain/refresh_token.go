package domain

import "time"

type RefreshToken struct {
	TokenID   string
	UserID    string
	ExpiresAt time.Time
}
