package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserProfile represents the user's profile information.
type UserProfile struct {
	ProfilePicture string    `bson:"profile_picture,omitempty" json:"profile_picture,omitempty"`
	Bio            string    `bson:"bio,omitempty" json:"bio,omitempty"`
	DateJoined     time.Time `bson:"date_joined,omitempty" json:"date_joined,omitempty"`
	LastLogin      time.Time `bson:"last_login,omitempty" json:"last_login,omitempty"`
}

// Preferences represents the user's application preferences.
type Preferences struct {
	Theme          string `bson:"theme,omitempty" json:"theme,omitempty"`
	Language       string `bson:"language,omitempty" json:"language,omitempty"`
	Notifications  bool   `bson:"notifications,omitempty" json:"notifications,omitempty"`
	TextSize       string `bson:"text_size,omitempty" json:"text_size,omitempty"`
	DefaultVoice   string `bson:"default_voice,omitempty" json:"default_voice,omitempty"`
	AvailableTTS   []string `bson:"available_tts,omitempty" json:"available_tts,omitempty"`
	AvailableCloned []string `bson:"available_cloned,omitempty" json:"available_cloned,omitempty"`
}

// SavedWord represents a word saved by the user.
type SavedWord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Word       string             `bson:"word,omitempty" json:"word,omitempty"`
	Definition string             `bson:"definition,omitempty" json:"definition,omitempty"`
	Context    string             `bson:"context,omitempty" json:"context,omitempty"`
	DateSaved  time.Time          `bson:"date_saved,omitempty" json:"date_saved,omitempty"`
	ReadingID  string             `bson:"reading_id,omitempty" json:"reading_id,omitempty"` // Or primitive.ObjectID if it references another collection
}

// Quiz represents quiz data within reading progress.
type Quiz struct {
	// Define fields for quiz data based on your original schema
	// Example:
	// Score int `bson:"score,omitempty" json:"score,omitempty"`
	// Completed bool `bson:"completed,omitempty" json:"completed,omitempty"`
}

// ReadingProgress represents the user's progress in a specific reading.
type ReadingProgress struct {
	LastPageRead    int    `bson:"last_page_read,omitempty" json:"last_page_read,omitempty"`
	CompletionRate  int    `bson:"completion_rate,omitempty" json:"completion_rate,omitempty"` // Assuming 0-100
	Quiz            Quiz   `bson:"quiz,omitempty" json:"quiz,omitempty"`
	CurrentSentence []int  `bson:"current_sentence,omitempty" json:"current_sentence,omitempty"`
}

// ClickedWord represents a word clicked by the user within a reading.
type ClickedWord struct {
	Text      string `bson:"text,omitempty" json:"text,omitempty"`
	Sentence  []int  `bson:"sentence,omitempty" json:"sentence,omitempty"` // [paragraph_index, sentence_index]
	Timestamp int64  `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

// UserReading represents a reading associated with a user.
type UserReading struct {
	Owner             bool               `bson:"owner" json:"owner"` // True if the user is the original owner
	ReadingInternalID int                `bson:"reading_internal_id" json:"reading_internal_id"`
	ReadingProgress   ReadingProgress    `bson:"reading_progress,omitempty" json:"reading_progress,omitempty"`
	ReadingID         primitive.ObjectID `bson:"reading_id" json:"reading_id"` // Reference to the Reading document
	ClickedWords      []ClickedWord      `bson:"clicked_words,omitempty" json:"clicked_words,omitempty"`
}

// User represents the main user schema.
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name              string             `bson:"name" json:"name"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"password"` // This will store the hashed password
	IsVerified        bool               `bson:"is_verified,omitempty" json:"is_verified,omitempty"`
	VerificationToken string             `bson:"verification_token,omitempty" json:"verification_token,omitempty"`
	ResetPasswordToken string            `bson:"reset_password_token,omitempty" json:"reset_password_token,omitempty"`
	ResetPasswordExpires time.Time      `bson:"reset_password_expires,omitempty" json:"reset_password_expires,omitempty"`
	Profile           UserProfile        `bson:"profile,omitempty" json:"profile,omitempty"`
	Preferences       Preferences        `bson:"preferences,omitempty" json:"preferences,omitempty"`
	Readings          []UserReading      `bson:"readings,omitempty" json:"readings,omitempty"`
	SavedWords        []SavedWord        `bson:"saved_words,omitempty" json:"saved_words,omitempty"`
	CreatedAt         time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt         time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
