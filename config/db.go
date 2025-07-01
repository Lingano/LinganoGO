package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// ConnectDB connects to the PostgreSQL database.
func ConnectDB() error {
	err := godotenv.Load() // Loads .env from the current directory
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("Successfully connected to PostgreSQL!")
	return nil
}

// GetDB returns the database connection pool.
func GetDB() *sql.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call ConnectDB first.")
	}
	return DB
}

// DisconnectDB closes the database connection.
func DisconnectDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
		log.Println("Disconnected from PostgreSQL.")
	}
}
