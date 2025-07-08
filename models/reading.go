package models

import (
	"LinganoGO/config"
	"fmt"

	"github.com/google/uuid"
)

type Reading struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	UserID   string `json:"userID"`
	Finished bool   `json:"finished"`
	Public   bool   `json:"public"`
}

func (r *Reading) Create() error {
	db := config.GetDB()
	r.ID = uuid.New().String()
	_, err := db.Exec("INSERT INTO readings (id, title, user_id, finished, public) VALUES ($1, $2, $3, $4, $5)", r.ID, r.Title, r.UserID, r.Finished, r.Public)
	return err
}

// GetAll fetches all readings from the database.
func GetAllReadings() ([]*Reading, error) {
	db := config.GetDB()
	rows, err := db.Query("SELECT id, title, user_id, finished, public FROM readings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*Reading
	for rows.Next() {
		var reading Reading
		if err := rows.Scan(&reading.ID, &reading.Title, &reading.UserID, &reading.Finished, &reading.Public); err != nil {
			return nil, err
		}
		readings = append(readings, &reading)
	}

	return readings, nil
}

// GetReadingByID fetches a reading from the database by its ID.
func GetReadingByID(id string) (*Reading, error) {
	db := config.GetDB()
	reading := &Reading{}

	query := `SELECT id, title, user_id, finished, public FROM readings WHERE id = $1`
	err := db.QueryRow(query, id).Scan(
		&reading.ID,
		&reading.Title,
		&reading.UserID,
		&reading.Finished,
		&reading.Public,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get reading: %w", err)
	}

	return reading, nil
}

// UpdatePublicStatus updates the public status of a reading.
func (r *Reading) UpdatePublicStatus(public bool) error {
	db := config.GetDB()
	_, err := db.Exec("UPDATE readings SET public = $1 WHERE id = $2", public, r.ID)
	if err != nil {
		return err
	}
	r.Public = public
	return nil
}

// UpdateReadingPublicStatus updates the public status of a reading.
func UpdateReadingPublicStatus(id string, public bool) error {
	db := config.GetDB()
	_, err := db.Exec("UPDATE readings SET public = $1 WHERE id = $2", public, id)
	if err != nil {
		return fmt.Errorf("failed to update reading public status: %w", err)
	}
	return nil
}

// GetPublicReadings fetches all public readings from the database.
func GetPublicReadings() ([]*Reading, error) {
	db := config.GetDB()
	rows, err := db.Query("SELECT id, title, user_id, finished, public FROM readings WHERE public = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*Reading
	for rows.Next() {
		var reading Reading
		if err := rows.Scan(&reading.ID, &reading.Title, &reading.UserID, &reading.Finished, &reading.Public); err != nil {
			return nil, err
		}
		readings = append(readings, &reading)
	}

	return readings, nil
}

// GetUserReadings fetches all readings for a specific user (both public and private).
func GetUserReadings(userID string) ([]*Reading, error) {
	db := config.GetDB()
	rows, err := db.Query("SELECT id, title, user_id, finished, public FROM readings WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*Reading
	for rows.Next() {
		var reading Reading
		if err := rows.Scan(&reading.ID, &reading.Title, &reading.UserID, &reading.Finished, &reading.Public); err != nil {
			return nil, err
		}
		readings = append(readings, &reading)
	}

	return readings, nil
}
