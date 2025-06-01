package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client

// ConnectDB connects to MongoDB and initializes MongoClient.
func ConnectDB() error {
	// Load .env file. In production, environment variables should be set directly.
	err := godotenv.Load() // Loads .env from the current directory
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	MongoClient = client
	log.Println("Successfully connected to MongoDB!")
	return nil
}

// GetCollection returns a handle to a specific collection.
func GetCollection(collectionName string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoDB client not initialized. Call ConnectDB first.")
	}
	// You might want to specify your database name here directly
	// or get it from an environment variable as well.
	databaseName := os.Getenv("DB_NAME")
	if databaseName == "" {
		databaseName = "lingano" // Default database name
		log.Printf("DB_NAME environment variable not set, using default: %s\n", databaseName)
	}
	return MongoClient.Database(databaseName).Collection(collectionName)
}

// DisconnectDB closes the MongoDB connection.
func DisconnectDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		log.Println("Disconnected from MongoDB.")
	}
}
