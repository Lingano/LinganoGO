package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"LinganoGO/ent"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var entClient *ent.Client

// ConnectEntDB initializes the Ent client with PostgreSQL
func ConnectEntDB() error {
	err := godotenv.Load(".env") // Loads .env from the current directory
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	
	client, err := ent.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed opening connection to postgres: %v", err)
	}
	
	entClient = client
	
	// Run the auto migration tool to create/update database schema
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Printf("Warning: failed creating schema resources: %v", err)
		// Don't return error here as tables might already exist
	}
	
	return nil
}

// GetEntClient returns the Ent client instance
func GetEntClient() *ent.Client {
	return entClient
}

// DisconnectEntDB closes the Ent client connection
func DisconnectEntDB() {
	if entClient != nil {
		entClient.Close()
	}
}
