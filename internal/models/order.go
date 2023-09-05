package models

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey" json:"-"`
	Number     string    `gorm:"not null;unique" json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UserID     uint      `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
}
