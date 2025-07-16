package services

import (
	"context"
	"fmt"

	"LinganoGO/config"
	"LinganoGO/ent"
	"LinganoGO/ent/reading"
	"LinganoGO/ent/user"

	"github.com/google/uuid"
)

// ReadingService provides methods for reading operations using Ent
type ReadingService struct {
	client *ent.Client
}

// NewReadingService creates a new ReadingService
func NewReadingService() *ReadingService {
	return &ReadingService{
		client: config.GetEntClient(),
	}
}

// CreateReading creates a new reading using Ent
func (s *ReadingService) CreateReading(ctx context.Context, title string, userID uuid.UUID, public bool) (*ent.Reading, error) {
	reading, err := s.client.Reading.
		Create().
		SetTitle(title).
		SetUserID(userID).
		SetPublic(public).
		SetFinished(false).
		Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create reading: %w", err)
	}
	
	return reading, nil
}

// GetReadingByID retrieves a reading by ID using Ent
func (s *ReadingService) GetReadingByID(ctx context.Context, id uuid.UUID) (*ent.Reading, error) {
	reading, err := s.client.Reading.
		Query().
		Where(reading.IDEQ(id)).
		First(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get reading: %w", err)
	}
	
	return reading, nil
}

// GetAllReadings retrieves all readings using Ent
func (s *ReadingService) GetAllReadings(ctx context.Context) ([]*ent.Reading, error) {
	readings, err := s.client.Reading.
		Query().
		All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get all readings: %w", err)
	}
	
	return readings, nil
}

// GetPublicReadings retrieves all public readings using Ent
func (s *ReadingService) GetPublicReadings(ctx context.Context) ([]*ent.Reading, error) {
	readings, err := s.client.Reading.
		Query().
		Where(reading.PublicEQ(true)).
		All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get public readings: %w", err)
	}
	
	return readings, nil
}

// GetReadingsByUser retrieves all readings for a specific user using Ent
func (s *ReadingService) GetReadingsByUser(ctx context.Context, userID uuid.UUID) ([]*ent.Reading, error) {
	readings, err := s.client.Reading.
		Query().
		Where(reading.HasUserWith(user.IDEQ(userID))).
		All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get readings by user: %w", err)
	}
	
	return readings, nil
}

// UpdateReading updates a reading using Ent
func (s *ReadingService) UpdateReading(ctx context.Context, id uuid.UUID, title string, finished, public bool) (*ent.Reading, error) {
	reading, err := s.client.Reading.
		UpdateOneID(id).
		SetTitle(title).
		SetFinished(finished).
		SetPublic(public).
		Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to update reading: %w", err)
	}
	
	return reading, nil
}

// MarkReadingAsFinished marks a reading as finished using Ent
func (s *ReadingService) MarkReadingAsFinished(ctx context.Context, id uuid.UUID) (*ent.Reading, error) {
	reading, err := s.client.Reading.
		UpdateOneID(id).
		SetFinished(true).
		Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to mark reading as finished: %w", err)
	}
	
	return reading, nil
}

// DeleteReading deletes a reading using Ent
func (s *ReadingService) DeleteReading(ctx context.Context, id uuid.UUID) error {
	err := s.client.Reading.
		DeleteOneID(id).
		Exec(ctx)
	
	if err != nil {
		return fmt.Errorf("failed to delete reading: %w", err)
	}
	
	return nil
}

// GetReadingWithUser retrieves a reading with its associated user using Ent edges
func (s *ReadingService) GetReadingWithUser(ctx context.Context, id uuid.UUID) (*ent.Reading, error) {
	reading, err := s.client.Reading.
		Query().
		Where(reading.IDEQ(id)).
		WithUser().
		First(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get reading with user: %w", err)
	}
	
	return reading, nil
}

// GetReadingsByUserID retrieves all readings for a specific user by user ID
func (s *ReadingService) GetReadingsByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.Reading, error) {
	readings, err := s.client.Reading.
		Query().
		Where(reading.UserIDEQ(userID)).
		All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get readings by user ID: %w", err)
	}
	
	return readings, nil
}

// DeleteReadingsByTitlePattern deletes readings matching a title pattern
func (s *ReadingService) DeleteReadingsByTitlePattern(ctx context.Context, pattern string) error {
	_, err := s.client.Reading.
		Delete().
		Where(reading.TitleContains(pattern)).
		Exec(ctx)
	
	if err != nil {
		return fmt.Errorf("failed to delete readings by pattern: %w", err)
	}
	
	return nil
}

// CountReadings returns the total number of readings
func (s *ReadingService) CountReadings(ctx context.Context) (int, error) {
	count, err := s.client.Reading.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count readings: %w", err)
	}
	
	return count, nil
}

// GetFinishedReadings retrieves all finished readings
func (s *ReadingService) GetFinishedReadings(ctx context.Context) ([]*ent.Reading, error) {
	readings, err := s.client.Reading.
		Query().
		Where(reading.FinishedEQ(true)).
		All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get finished readings: %w", err)
	}
	
	return readings, nil
}
