package services

import (
	"LinganoGO/config"
	"LinganoGO/ent"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type FlashcardService struct {
	client *ent.Client
}

func NewFlashcardService() *FlashcardService {
	return &FlashcardService{
		client: config.GetEntClient(),
	}
}

func (s *FlashcardService) UpdateFlashcard(ctx context.Context, id string, question string, answer string) (*ent.Flashcard, error) {
	flashcardUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid flashcard ID: %w", err)
	}
	flashcard, err := s.client.Flashcard.
		UpdateOneID(flashcardUUID).
		SetQuestion(question).
		SetAnswer(answer).
		Save(ctx)
		
	if err != nil {
		return nil, fmt.Errorf("failed to update flashcard: %w", err)
	}
	return flashcard, nil
}
