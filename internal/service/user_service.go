package service

import (
	"context"
	"fmt"
	"time"

	"jira-clone-be/internal/repository"
	"jira-clone-be/pkg/logger"
	"jira-clone-be/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, logger *logger.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *types.CreateUserRequest) (*types.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password (in a real implementation, you'd store this)
	_, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to hash password")
	}

	// Create user
	user := &types.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		AvatarURL: s.generateAvatarURL(req.Name),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IssueIDs:  []string{},
	}

	// Store hashed password (in a real implementation, you'd store this in a separate table)
	// For now, we'll just create the user without storing the password
	// In a real implementation, you'd have a separate UserAuth table

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User created successfully", map[string]interface{}{
		"userID": user.ID,
		"email":  user.Email,
	})

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*types.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get user", map[string]interface{}{
			"error":  err.Error(),
			"userID": id,
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get user by email", map[string]interface{}{
			"error": err.Error(),
			"email": email,
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUsers retrieves all users
func (s *UserService) GetUsers(ctx context.Context) ([]*types.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get users", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id string, req *types.UpdateUserRequest) (*types.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get user for update", map[string]interface{}{
			"error":  err.Error(),
			"userID": id,
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now()

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user", map[string]interface{}{
			"error":  err.Error(),
			"userID": id,
		})
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info("User updated successfully", map[string]interface{}{
		"userID": user.ID,
		"email":  user.Email,
	})

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get user for deletion", map[string]interface{}{
			"error":  err.Error(),
			"userID": id,
		})
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete user", map[string]interface{}{
			"error":  err.Error(),
			"userID": id,
		})
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("User deleted successfully", map[string]interface{}{
		"userID": id,
	})

	return nil
}

// ValidatePassword validates a password against a hashed password
func (s *UserService) ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// generateAvatarURL generates an avatar URL based on the user's name
func (s *UserService) generateAvatarURL(name string) string {
	// In a real implementation, you'd generate a proper avatar URL
	// For now, we'll use a placeholder
	return fmt.Sprintf("/assets/avatars/%s.png", name)
}
