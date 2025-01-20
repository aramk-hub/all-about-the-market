package database

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"all-about-the-market/backend/models" // Adjust with your actual project path
)

var DB *gorm.DB

func ConnectDB() {
	// Define the connection string (DSN)
	dsn := "host=localhost user=postgres password=kridershot dbname=portfolio_tracker port=5432 sslmode=disable"
	var err error

	// Open the database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Automatically migrate schema (create tables, etc.)
	err = DB.AutoMigrate(
		&models.Portfolio{}, 
		&models.Stock{}, 
		&models.Portfolio{})

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully!")
}
