package models

import "time"

type StoredCode struct {
	Code      string    `json:"code"`
	UserID    string    `json:"user_id"`
	UserEmail string    `json:"user_email"`
	CreatedAt time.Time `json:"created_at"` // The start of the UTC time window
}
