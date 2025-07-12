package main

import (
	"log"

	"LinganoGO/config"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
	
)

func main() {
	// Connect to PostgreSQL
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer config.DisconnectDB()

	db := config.GetDB()
	
	// Ensure uuid-ossp extension is created
	if _, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		log.Fatalf("Failed to create uuid-ossp extension: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	// Get the path to the migrations directory
	migrationsDir := "./migrations"

	// Run migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}
