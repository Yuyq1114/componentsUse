package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Product   string
	Price     float64
	CreatedAt time.Time
}
