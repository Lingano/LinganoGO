package main

import (
	"context"
	"log"

	"LinganoGO/config"
	"LinganoGO/services"
)

func main() {
	log.Println("🧪 Testing Ent Setup")
	log.Println("====================")

	// Connect to database
	if err := config.ConnectEntDB(); err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer config.DisconnectEntDB()

	ctx := context.Background()
	userService := services.NewUserService()
	readingService := services.NewReadingService()

	// Test 1: Create a user
	log.Println("\n1. Creating a user...")
	user, err := userService.CreateUser(ctx, "John Doe", "john@test.com", "password123")
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}
	log.Printf("✅ Created: %s (ID: %s)", user.Name, user.ID)

	// Test 2: Create a reading
	log.Println("\n2. Creating a reading...")
	reading, err := readingService.CreateReading(ctx, "Test Book", user.ID, true)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}
	log.Printf("✅ Created: %s (Public: %t)", reading.Title, reading.Public)

	// Test 3: Get all users
	log.Println("\n3. Getting all users...")
	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}
	log.Printf("✅ Found %d users:", len(users))
	for i, u := range users {
		log.Printf("   %d. %s (%s)", i+1, u.Name, u.Email)
	}

	// Test 4: Get user with readings
	log.Println("\n4. Getting user with readings...")
	userWithReadings, err := userService.GetUserWithReadings(ctx, user.ID)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}
	log.Printf("✅ %s has %d readings:", userWithReadings.Name, len(userWithReadings.Edges.Readings))
	for i, r := range userWithReadings.Edges.Readings {
		log.Printf("   %d. %s", i+1, r.Title)
	}

	// Test 5: Get public readings
	log.Println("\n5. Getting public readings...")
	publicReadings, err := readingService.GetPublicReadings(ctx)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}
	log.Printf("✅ Found %d public readings:", len(publicReadings))
	for i, r := range publicReadings {
		log.Printf("   %d. %s", i+1, r.Title)
	}

	log.Println("\n🎉 All tests passed! Ent is working correctly.")
}
