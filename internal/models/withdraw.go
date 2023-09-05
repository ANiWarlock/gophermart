package models

import (
	"time"
)

type Withdraw struct {
	ID          uint      `gorm:"primaryKey" json:"-"`
	UserID      uint      `gorm:"not null" json:"-"`
	Order       string    `gorm:"size:255;not null;unique" json:"order"`
	Sum         float64   `gorm:"not null;" json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
