package models

type Balance struct {
	ID        uint    `gorm:"primaryKey" json:"-"`
	UserID    uint    `gorm:"not null;unique;" json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
