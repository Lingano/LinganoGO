package services

import (
	"LinganoGO/config"
	"LinganoGO/ent"
)

type PostService struct {
	client *ent.Client
}

func NewPostService() *PostService {
	return &PostService{
		client: config.GetEntClient(),
	}
}

