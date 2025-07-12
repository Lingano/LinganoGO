package main

import (
	"context"
	"log"
	"strings"

	"LinganoGO/config"
	"LinganoGO/ent"
	"LinganoGO/ent/reading"
	"LinganoGO/ent/user"
	"LinganoGO/services"
)

func main() {
	// Connect to database
	log.Println("🔌 Connecting to database...")
	if err := config.ConnectEntDB(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer config.DisconnectEntDB()

	ctx := context.Background()
	client := config.GetEntClient()

	// Initialize services
	userService := services.NewUserService()
	readingService := services.NewReadingService()

	log.Println("🧪 Starting Ent Query Tests...")
	log.Println(strings.Repeat("=", 50))

	// Test 1: Basic User Operations
	testBasicUserOperations(ctx, userService)

	// Test 2: Basic Reading Operations
	testBasicReadingOperations(ctx, readingService)

	// Test 3: Relationship Queries
	testRelationshipQueries(ctx, client)

	// Test 4: Advanced Queries
	testAdvancedQueries(ctx, client)

	// Test 5: Aggregations
	testAggregations(ctx, client)

	// Test 6: Filtering and Sorting
	testFilteringAndSorting(ctx, client)

	log.Println("🎉 All tests completed!")
}

func testBasicUserOperations(ctx context.Context, userService *services.UserService) {
	log.Println("\n📋 Test 1: Basic User Operations")
	log.Println(strings.Repeat("-", 30))

	// Create users
	log.Println("Creating users...")
	user1, err := userService.CreateUser(ctx, "Alice Johnson", "alice@example.com", "hashedpassword1")
	if err != nil {
		log.Printf("❌ Error creating user1: %v", err)
		return
	}
	log.Printf("✅ Created user: %s (ID: %s)", user1.Name, user1.ID)

	user2, err := userService.CreateUser(ctx, "Bob Smith", "bob@example.com", "hashedpassword2")
	if err != nil {
		log.Printf("❌ Error creating user2: %v", err)
		return
	}
	log.Printf("✅ Created user: %s (ID: %s)", user2.Name, user2.ID)

	// Get user by ID
	log.Println("\nFetching user by ID...")
	fetchedUser, err := userService.GetUserByID(ctx, user1.ID)
	if err != nil {
		log.Printf("❌ Error fetching user: %v", err)
	} else {
		log.Printf("✅ Fetched user: %s (%s)", fetchedUser.Name, fetchedUser.Email)
	}

	// Get user by email
	log.Println("\nFetching user by email...")
	userByEmail, err := userService.GetUserByEmail(ctx, "bob@example.com")
	if err != nil {
		log.Printf("❌ Error fetching user by email: %v", err)
	} else {
		log.Printf("✅ Found user by email: %s", userByEmail.Name)
	}

	// Get all users
	log.Println("\nFetching all users...")
	allUsers, err := userService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("❌ Error fetching all users: %v", err)
	} else {
		log.Printf("✅ Total users: %d", len(allUsers))
		for i, u := range allUsers {
			log.Printf("   %d. %s (%s)", i+1, u.Name, u.Email)
		}
	}
}

func testBasicReadingOperations(ctx context.Context, readingService *services.ReadingService) {
	log.Println("\n📚 Test 2: Basic Reading Operations")
	log.Println(strings.Repeat("-", 30))

	// First, get a user to associate readings with
	client := config.GetEntClient()
	users, err := client.User.Query().Limit(1).All(ctx)
	if err != nil || len(users) == 0 {
		log.Println("❌ No users found, skipping reading tests")
		return
	}
	testUser := users[0]

	// Create readings
	log.Println("Creating readings...")
	reading1, err := readingService.CreateReading(ctx, "The Go Programming Language", testUser.ID, true)
	if err != nil {
		log.Printf("❌ Error creating reading1: %v", err)
		return
	}
	log.Printf("✅ Created reading: %s (Public: %t)", reading1.Title, reading1.Public)

	reading2, err := readingService.CreateReading(ctx, "Advanced Go Patterns", testUser.ID, false)
	if err != nil {
		log.Printf("❌ Error creating reading2: %v", err)
		return
	}
	log.Printf("✅ Created reading: %s (Public: %t)", reading2.Title, reading2.Public)

	// Mark one as finished
	log.Println("\nMarking reading as finished...")
	finishedReading, err := readingService.MarkReadingAsFinished(ctx, reading1.ID)
	if err != nil {
		log.Printf("❌ Error marking as finished: %v", err)
	} else {
		log.Printf("✅ Marked as finished: %s (Finished: %t)", finishedReading.Title, finishedReading.Finished)
	}

	// Get public readings
	log.Println("\nFetching public readings...")
	publicReadings, err := readingService.GetPublicReadings(ctx)
	if err != nil {
		log.Printf("❌ Error fetching public readings: %v", err)
	} else {
		log.Printf("✅ Public readings: %d", len(publicReadings))
		for i, r := range publicReadings {
			log.Printf("   %d. %s (Finished: %t)", i+1, r.Title, r.Finished)
		}
	}
}

func testRelationshipQueries(ctx context.Context, client *ent.Client) {
	log.Println("\n🔗 Test 3: Relationship Queries")
	log.Println(strings.Repeat("-", 30))

	// Get user with all their readings
	log.Println("Fetching users with their readings...")
	usersWithReadings, err := client.User.
		Query().
		WithReadings().
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error fetching users with readings: %v", err)
		return
	}

	for _, u := range usersWithReadings {
		log.Printf("✅ User: %s has %d readings", u.Name, len(u.Edges.Readings))
		for i, r := range u.Edges.Readings {
			log.Printf("   %d. %s (Finished: %t, Public: %t)", i+1, r.Title, r.Finished, r.Public)
		}
	}

	// Get readings with their users
	log.Println("\nFetching readings with their users...")
	readingsWithUsers, err := client.Reading.
		Query().
		WithUser().
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error fetching readings with users: %v", err)
		return
	}

	for _, r := range readingsWithUsers {
		log.Printf("✅ Reading: '%s' by %s", r.Title, r.Edges.User.Name)
	}
}

func testAdvancedQueries(ctx context.Context, client *ent.Client) {
	log.Println("\n🚀 Test 4: Advanced Queries")
	log.Println(strings.Repeat("-", 30))

	// Find users who have finished readings
	log.Println("Finding users with finished readings...")
	usersWithFinished, err := client.User.
		Query().
		Where(user.HasReadingsWith(reading.FinishedEQ(true))).
		WithReadings(func(q *ent.ReadingQuery) {
			q.Where(reading.FinishedEQ(true))
		}).
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		log.Printf("✅ Users with finished readings: %d", len(usersWithFinished))
		for _, u := range usersWithFinished {
			log.Printf("   - %s has %d finished readings", u.Name, len(u.Edges.Readings))
		}
	}

	// Find public readings by specific users
	log.Println("\nFinding public readings...")
	publicReadings, err := client.Reading.
		Query().
		Where(
			reading.PublicEQ(true),
			reading.HasUserWith(user.NameContains("Alice")),
		).
		WithUser().
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		log.Printf("✅ Public readings by Alice: %d", len(publicReadings))
		for _, r := range publicReadings {
			log.Printf("   - '%s' by %s", r.Title, r.Edges.User.Name)
		}
	}

	// Complex query: Unfinished public readings
	log.Println("\nFinding unfinished public readings...")
	unfinishedPublic, err := client.Reading.
		Query().
		Where(
			reading.PublicEQ(true),
			reading.FinishedEQ(false),
		).
		WithUser().
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		log.Printf("✅ Unfinished public readings: %d", len(unfinishedPublic))
		for _, r := range unfinishedPublic {
			log.Printf("   - '%s' by %s", r.Title, r.Edges.User.Name)
		}
	}
}

func testAggregations(ctx context.Context, client *ent.Client) {
	log.Println("\n📊 Test 5: Aggregations")
	log.Println(strings.Repeat("-", 30))

	// Count total users
	userCount, err := client.User.Query().Count(ctx)
	if err != nil {
		log.Printf("❌ Error counting users: %v", err)
	} else {
		log.Printf("✅ Total users: %d", userCount)
	}

	// Count total readings
	readingCount, err := client.Reading.Query().Count(ctx)
	if err != nil {
		log.Printf("❌ Error counting readings: %v", err)
	} else {
		log.Printf("✅ Total readings: %d", readingCount)
	}

	// Count finished readings
	finishedCount, err := client.Reading.Query().Where(reading.FinishedEQ(true)).Count(ctx)
	if err != nil {
		log.Printf("❌ Error counting finished: %v", err)
	} else {
		log.Printf("✅ Finished readings: %d", finishedCount)
	}

	// Count public readings
	publicCount, err := client.Reading.Query().Where(reading.PublicEQ(true)).Count(ctx)
	if err != nil {
		log.Printf("❌ Error counting public: %v", err)
	} else {
		log.Printf("✅ Public readings: %d", publicCount)
	}

	// Count readings per user
	log.Println("\nReadings per user:")
	users, err := client.User.Query().WithReadings().All(ctx)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		for _, u := range users {
			log.Printf("   - %s: %d readings", u.Name, len(u.Edges.Readings))
		}
	}
}

func testFilteringAndSorting(ctx context.Context, client *ent.Client) {
	log.Println("\n🔍 Test 6: Filtering and Sorting")
	log.Println(strings.Repeat("-", 30))

	// Get users sorted by name
	log.Println("Users sorted by name:")
	usersByName, err := client.User.
		Query().
		Order(ent.Asc(user.FieldName)).
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		for i, u := range usersByName {
			log.Printf("   %d. %s (%s)", i+1, u.Name, u.Email)
		}
	}

	// Get readings sorted by title
	log.Println("\nReadings sorted by title:")
	readingsByTitle, err := client.Reading.
		Query().
		Order(ent.Asc(reading.FieldTitle)).
		WithUser().
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		for i, r := range readingsByTitle {
			log.Printf("   %d. '%s' by %s", i+1, r.Title, r.Edges.User.Name)
		}
	}

	// Limit and offset
	log.Println("\nFirst 2 users (pagination):")
	firstUsers, err := client.User.
		Query().
		Order(ent.Asc(user.FieldName)).
		Limit(2).
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		for i, u := range firstUsers {
			log.Printf("   %d. %s", i+1, u.Name)
		}
	}

	// Search by name pattern
	log.Println("\nUsers with 'A' in name:")
	usersWithA, err := client.User.
		Query().
		Where(user.NameContains("A")).
		All(ctx)
	
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		for _, u := range usersWithA {
			log.Printf("   - %s", u.Name)
		}
	}
}
