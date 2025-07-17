package services

import (
	"LinganoGO/config"
	"LinganoGO/ent"
)

type FlashcardService struct {
	client *ent.Client
}

func NewFlashcardService() *FlashcardService {
	return &FlashcardService{
		client: config.GetEntClient(),
	}
}