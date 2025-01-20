package models

import "time"

type Portfolio struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null"`
	Name        string     `gorm:"not null"`
	Description string     `gorm:"size:255"`
	Stocks      []Stock    `gorm:"foreignKey:PortfolioID"` // One-to-many relationship with stocks
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
