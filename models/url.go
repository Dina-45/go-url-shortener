package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	OriginalURL string `json:"original_url" gorm:"not null"`
	ShortCode   string `json:"short_code" gorm:"unique;not null"`
	Clicks      int    `json:"clicks"`
}
