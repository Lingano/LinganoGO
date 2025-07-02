package models

import (
	"database/sql"
	"fmt"
	"time"

	"LinganoGO/config"
)

// User represents the user schema in PostgreSQL.
type User struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	Password             string    `json:"password"` // Hashed password
	IsVerified           bool      `json:"is_verified"`
	VerificationToken    string    `json:"verification_token"`
	ResetPasswordToken   string    `json:"reset_password_token"`
	ResetPasswordExpires sql.NullTime `json:"reset_password_expires"`
	Profile              []byte    `json:"profile"`    // Stored as JSONB
	Preferences          []byte    `json:"preferences"` // Stored as JSONB
	Readings             []byte    `json:"readings"`   // Stored as JSONB
	SavedWords           []byte    `json:"saved_words"` // Stored as JSONB
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// GetUserByID fetches a user from the database by their ID.
func GetUserByID(id string) (*User, error) {
	db := config.GetDB()
	user := &User{}

	query := `SELECT id, name, email, password, is_verified, verification_token, reset_password_token, reset_password_expires, profile, preferences, readings, saved_words, created_at, updated_at FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.VerificationToken,
		&user.ResetPasswordToken,
		&user.ResetPasswordExpires,
		&user.Profile,
		&user.Preferences,
		&user.Readings,
		&user.SavedWords,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil
}

// Create inserts a new user into the database.
func (u *User) Create() error {
	db := config.GetDB()
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := db.QueryRow(query, u.Name, u.Email, u.Password).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
