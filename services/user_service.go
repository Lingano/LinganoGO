package services

import (
	"context"
	"fmt"

	"LinganoGO/config"
	"LinganoGO/ent"
	"LinganoGO/ent/user"

	"github.com/google/uuid"
)

// UserService provides methods for user operations using Ent
type UserService struct {
	client *ent.Client
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		client: config.GetEntClient(),
	}
}

// CreateUser creates a new user using Ent
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*ent.User, error) {
	user, err := s.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPassword(password).
		SetRole(user.RoleUSER).
		Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID using Ent
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	user, err := s.client.User.
		Query().
		Where(user.IDEQ(id)).
		First(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

// GetUserByEmail retrieves a user by email using Ent
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	user, err := s.client.User.
		Query().
		Where(user.EmailEQ(email)).
		First(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

// GetAllUsers retrieves all users using Ent
func (s *UserService) GetAllUsers(ctx context.Context) ([]*ent.User, error) {
	users, err := s.client.User.
		Query().
		All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	
	return users, nil
}

// UpdateUser updates a user using Ent
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (*ent.User, error) {
	user, err := s.client.User.
		UpdateOneID(id).
		SetName(name).
		SetEmail(email).
		Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	
	return user, nil
}

// DeleteUser deletes a user using Ent
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.client.User.
		DeleteOneID(id).
		Exec(ctx)
	
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}

// GetUserWithReadings retrieves a user with their readings using Ent edges
func (s *UserService) GetUserWithReadings(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	user, err := s.client.User.
		Query().
		Where(user.IDEQ(id)).
		WithReadings().
		First(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user with readings: %w", err)
	}
	
	return user, nil
}

// ListUsers retrieves users with pagination using Ent
func (s *UserService) ListUsers(ctx context.Context, offset, limit int) ([]*ent.User, error) {
	users, err := s.client.User.
		Query().
		Offset(offset).
		Limit(limit).
		All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	
	return users, nil
}

// DeleteUsersByEmailPattern deletes users matching an email pattern
func (s *UserService) DeleteUsersByEmailPattern(ctx context.Context, pattern string) error {
	_, err := s.client.User.
		Delete().
		Where(user.EmailContains(pattern)).
		Exec(ctx)
	
	if err != nil {
		return fmt.Errorf("failed to delete users by pattern: %w", err)
	}
	
	return nil
}

// CountUsers returns the total number of users
func (s *UserService) CountUsers(ctx context.Context) (int, error) {
	count, err := s.client.User.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	return count, nil
}
