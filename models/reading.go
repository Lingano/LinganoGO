package models

import (
	"LinganoGO/config"

	"github.com/google/uuid"
)

type Reading struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	UserID   string `json:"userID"`
	Finished bool   `json:"finished"`
}

func (r *Reading) Create() error {
	db := config.GetDB()
	r.ID = uuid.New().String()
	_, err := db.Exec("INSERT INTO readings (id, title, user_id, finished) VALUES ($1, $2, $3, $4)", r.ID, r.Title, r.UserID, r.Finished)
	return err
}

// GetAll fetches all readings from the database.
func GetAllReadings() ([]*Reading, error) {
	db := config.GetDB()
	rows, err := db.Query("SELECT id, title, user_id, finished FROM readings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*Reading
	for rows.Next() {
		var reading Reading
		if err := rows.Scan(&reading.ID, &reading.Title, &reading.UserID, &reading.Finished); err != nil {
			return nil, err
		}
		readings = append(readings, &reading)
	}

	return readings, nil
}
