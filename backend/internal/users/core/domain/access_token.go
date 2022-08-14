package domain

import "time"

type AccessToken struct {
	UserID    string
	ExpiresAt time.Time
}
