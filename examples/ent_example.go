package main

import (
	"context"
	"log"

	"LinganoGO/config"
	"LinganoGO/services"
)

func main() {
	// Connect to PostgreSQL with Ent
	if err := config.ConnectEntDB(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL with Ent: %v", err)
	}
	defer config.DisconnectEntDB()

	ctx := context.Background()
	
	// Initialize services
	userService := services.NewUserService()
	readingService := services.NewReadingService()

	// Example: Create a user
	log.Println("Creating a user...")
	user, err := userService.CreateUser(ctx, "John Doe", "john@example.com", "hashedpassword123")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	} else {
		log.Printf("Created user: %s (ID: %s)", user.Name, user.ID)
	}

	// Example: Create a reading for the user
	if user != nil {
		log.Println("Creating a reading...")
		reading, err := readingService.CreateReading(ctx, "My First Reading", user.ID, true)
		if err != nil {
			log.Printf("Error creating reading: %v", err)
		} else {
			log.Printf("Created reading: %s (ID: %s)", reading.Title, reading.ID)
		}

		// Example: Get user with their readings
		log.Println("Fetching user with readings...")
		userWithReadings, err := userService.GetUserWithReadings(ctx, user.ID)
		if err != nil {
			log.Printf("Error fetching user with readings: %v", err)
		} else {
			log.Printf("User %s has %d readings", userWithReadings.Name, len(userWithReadings.Edges.Readings))
			for _, r := range userWithReadings.Edges.Readings {
				log.Printf("  - %s (finished: %t, public: %t)", r.Title, r.Finished, r.Public)
			}
		}
	}

	// Example: Get all public readings
	log.Println("Fetching all public readings...")
	publicReadings, err := readingService.GetPublicReadings(ctx)
	if err != nil {
		log.Printf("Error fetching public readings: %v", err)
	} else {
		log.Printf("Found %d public readings", len(publicReadings))
		for _, r := range publicReadings {
			log.Printf("  - %s by user %s", r.Title, r.UserID)
		}
	}

	log.Println("Ent example completed!")
}
