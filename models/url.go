package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code" gorm:"uniqueIndex"`
	IsAlive     *bool      `json:"is_alive"`
	Clicks      int        `json:"clicks"`
	UserID      uint       `json:"user_id"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
