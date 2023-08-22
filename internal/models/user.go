package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Login    string `gorm:"size:255;not null;unique" json:"login"`
	Password string `gorm:"size:255;not null;unique" json:"password"`
}
