package models

import "time"

// Stock represents an individual stock in a portfolio
type Stock struct {
	ID           uint      `gorm:"primaryKey"`
	PortfolioID  uint      `gorm:"not null"`                // Foreign key linking to Portfolio
	Symbol       string    `gorm:"not null;unique"`         // Stock ticker symbol (e.g., AAPL, MSFT)
	Shares       float64   `gorm:"not null"`                // Number of shares held
	AveragePrice float64   `gorm:"not null"`                // Average purchase price per share
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
