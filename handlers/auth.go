package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"LinganoGO/config"
	"LinganoGO/models"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo" // Added mongo import
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest defines the expected structure for a registration request.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body  handlers.RegisterRequest  true  "User info"
// @Success      201  {object}  models.UserResponse
// @Failure      400  {object}  models.ErrorResponse
// @Router       /api/auth/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var regReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	if regReq.Name == "" || regReq.Email == "" || regReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Name, email, and password are required"})
		return
	}

	// Get the users collection
	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"email": regReq.Email}).Decode(&existingUser)
	if err == nil {
		// User found
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "User with this email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// Another error occurred
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error checking for existing user"})
		return
	}
	// If mongo.ErrNoDocuments, then user does not exist, proceed.

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regReq.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to hash password"})
		return
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Name:      regReq.Name,
		Email:     regReq.Email,
		Password:  string(hashedPassword),
		IsVerified: false, // Default to not verified
		Profile:   models.UserProfile{DateJoined: time.Now()},
		Preferences: models.Preferences{}, // Initialize with default preferences if any
		Readings:  []models.UserReading{},
		SavedWords:[]models.SavedWord{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = usersCollection.InsertOne(ctx, newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "User registered successfully", "userId": newUser.ID.Hex()})
}

// LoginRequest defines the expected structure for a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims defines the JWT claims structure.
type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// LoginUser handles user login.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	if loginReq.Email == "" || loginReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email and password are required"})
		return
	}

	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var foundUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"email": loginReq.Email}).Decode(&foundUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error finding user"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginReq.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Generate JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Log this error, but don't expose it to the client directly for security reasons
		// log.Println("JWT_SECRET not set in environment variables")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not process login"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &Claims{
		UserID: foundUser.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "userId": foundUser.ID.Hex(), "name": foundUser.Name})
}
