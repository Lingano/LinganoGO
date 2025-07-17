package graph

import (
	"LinganoGO/ent"
	"LinganoGO/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	readings    []*ent.Reading
	userService *services.UserService
}

// NewResolver creates a new resolver with initialized services
func NewResolver() *Resolver {
	return &Resolver{
		userService: services.NewUserService(),
	}
}
