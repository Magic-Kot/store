package models

import "time"

type Request struct {
	UserId      int           `json:"user_id"`
	URL         string        `json:"url"`
	CustomShort string        `json:"short_link"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short_link"`
	Expiry      time.Duration `json:"expiry"`
}
