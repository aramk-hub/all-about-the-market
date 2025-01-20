package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`  // Username field (optional, could be unique)
	Email     string    `gorm:"not null;unique" json:"email"`  // Email field (unique identifier)
	Password  string    `json:"password"`  // Password for user registration (hashed in the database)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
