package services

import (
	"LinganoGO/config"
	"LinganoGO/ent"
	"LinganoGO/graph/model"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PostService struct {
	client *ent.Client
}

func NewPostService() *PostService {
	return &PostService{
		client: config.GetEntClient(),
	}
}


func (s *PostService) CreatePost(ctx context.Context, input model.NewPost) (*ent.Post, error) {
	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	var draft bool
	if input.Draft != nil {
		draft = *input.Draft
	} else {
		draft = false
	}
	post, err := s.client.Post.
		Create().
		SetBody(input.Body).
		SetDraft(draft).
		SetUserID(userUUID).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}