package models

import "time"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	RefreshToken string
	ExpiresAt    time.Time
}
