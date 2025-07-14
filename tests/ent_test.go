package tests

import (
	"context"
	"testing"

	"LinganoGO/config"
	"LinganoGO/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntIntegration(t *testing.T) {
	// Setup: Connect to database
	err := config.ConnectEntDB()
	require.NoError(t, err, "Failed to connect to Ent database")
	defer config.DisconnectEntDB()

	ctx := context.Background()
	userService := services.NewUserService()
	readingService := services.NewReadingService()

	t.Run("CreateUser", func(t *testing.T) {
		t.Log("Creating a user...")

		user, err := userService.CreateUser(ctx, "John Doe", "john1@test.com", "password123")
		require.NoError(t, err, "Failed to create user")

		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@test.com", user.Email)
		assert.NotEmpty(t, user.ID)

		t.Logf("✅ Created: %s (ID: %s)", user.Name, user.ID)
	})

	t.Run("CreateReading", func(t *testing.T) {
		t.Log("Creating a reading...")

		// First create a user for the reading
		user, err := userService.CreateUser(ctx, "Jane Doe", "jane1@test.com", "password123")
		require.NoError(t, err, "Failed to create user for reading test")

		reading, err := readingService.CreateReading(ctx, "Test Book", user.ID, true)
		require.NoError(t, err, "Failed to create reading")

		assert.Equal(t, "Test Book", reading.Title)
		assert.Equal(t, user.ID, reading.UserID)
		assert.True(t, reading.Public)
		assert.NotEmpty(t, reading.ID)

		t.Logf("✅ Created: %s (Public: %t)", reading.Title, reading.Public)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		t.Log("Getting all users...")

		users, err := userService.GetAllUsers(ctx)
		require.NoError(t, err, "Failed to get all users")

		assert.NotEmpty(t, users, "Expected at least one user")

		t.Logf("✅ Found %d users:", len(users))
		for i, u := range users {
			t.Logf("   %d. %s (%s)", i+1, u.Name, u.Email)
		}
	})

	t.Run("GetUserWithReadings", func(t *testing.T) {
		t.Log("Getting user with readings...")

		// Create a user and reading for this test
		user, err := userService.CreateUser(ctx, "Alice Smith", "alice1@test.com", "password123")
		require.NoError(t, err, "Failed to create user")

		reading, err := readingService.CreateReading(ctx, "Alice's Book", user.ID, false)
		require.NoError(t, err, "Failed to create reading")

		// Test getting user with readings
		userWithReadings, err := userService.GetUserWithReadings(ctx, user.ID)
		require.NoError(t, err, "Failed to get user with readings")

		assert.Equal(t, user.Name, userWithReadings.Name)
		assert.NotEmpty(t, userWithReadings.Edges.Readings, "Expected user to have readings")

		// Verify the reading is included
		found := false
		for _, r := range userWithReadings.Edges.Readings {
			if r.ID == reading.ID {
				found = true
				assert.Equal(t, "Alice's Book", r.Title)
				break
			}
		}
		assert.True(t, found, "Expected to find the created reading in user's readings")

		t.Logf("✅ %s has %d readings:", userWithReadings.Name, len(userWithReadings.Edges.Readings))
		for i, r := range userWithReadings.Edges.Readings {
			t.Logf("   %d. %s", i+1, r.Title)
		}
	})

	t.Run("GetPublicReadings", func(t *testing.T) {
		t.Log("Getting public readings...")

		// Create a user and public reading for this test
		user, err := userService.CreateUser(ctx, "Bob Wilson", "bob1@test.com", "password123")
		require.NoError(t, err, "Failed to create user")

		publicReading, err := readingService.CreateReading(ctx, "Public Book", user.ID, true)
		require.NoError(t, err, "Failed to create public reading")

		privateReading, err := readingService.CreateReading(ctx, "Private Book", user.ID, false)
		require.NoError(t, err, "Failed to create private reading")

		// Test getting public readings
		publicReadings, err := readingService.GetPublicReadings(ctx)
		require.NoError(t, err, "Failed to get public readings")

		assert.NotEmpty(t, publicReadings, "Expected at least one public reading")

		// Verify only public readings are returned
		foundPublic := false
		foundPrivate := false
		for _, r := range publicReadings {
			assert.True(t, r.Public, "All returned readings should be public")
			if r.ID == publicReading.ID {
				foundPublic = true
			}
			if r.ID == privateReading.ID {
				foundPrivate = true
			}
		}

		assert.True(t, foundPublic, "Expected to find the public reading")
		assert.False(t, foundPrivate, "Private reading should not be in public readings list")

		t.Logf("✅ Found %d public readings:", len(publicReadings))
		for i, r := range publicReadings {
			t.Logf("   %d. %s", i+1, r.Title)
		}
	})
}

// TestEntSetup tests basic Ent database connectivity
func TestEntSetup(t *testing.T) {
	t.Log("Testing Ent database setup...")

	err := config.ConnectEntDB()
	require.NoError(t, err, "Failed to connect to Ent database")
	defer config.DisconnectEntDB()

	client := config.GetEntClient()
	require.NotNil(t, client, "Ent client should not be nil")

	t.Log("✅ Ent database setup successful")
}

// BenchmarkUserCreation benchmarks user creation performance
func BenchmarkUserCreation(b *testing.B) {
	err := config.ConnectEntDB()
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.DisconnectEntDB()

	ctx := context.Background()
	userService := services.NewUserService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.CreateUser(ctx, "Benchmark User", "bench@test.com", "password")
	}
}

// BenchmarkReadingCreation benchmarks reading creation performance
func BenchmarkReadingCreation(b *testing.B) {
	err := config.ConnectEntDB()
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.DisconnectEntDB()

	ctx := context.Background()
	userService := services.NewUserService()
	readingService := services.NewReadingService()

	// Create a user for the readings
	user, err := userService.CreateUser(ctx, "Test User", "test@bench.com", "password")
	if err != nil {
		b.Fatalf("Failed to create user: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readingService.CreateReading(ctx, "Benchmark Book", user.ID, true)
	}
}
