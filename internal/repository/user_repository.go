package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"jira-clone-be/pkg/logger"
	"jira-clone-be/types"

	"github.com/google/uuid"
)

// userRepository implements UserRepository interface using JSON file storage
type userRepository struct {
	filePath string
	users    map[string]*types.User
	mutex    sync.RWMutex
	logger   *logger.Logger
}

// NewUserRepository creates a new user repository instance
func NewUserRepository() UserRepository {
	repo := &userRepository{
		filePath: "data/users.json",
		users:    make(map[string]*types.User),
		logger:   logger.New("info"),
	}

	// Load existing data
	if err := repo.loadData(); err != nil {
		repo.logger.Warn("Failed to load user data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return repo
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *types.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user with email already exists
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Initialize empty slice for IssueIDs
	if user.IssueIDs == nil {
		user.IssueIDs = []string{}
	}

	r.users[user.ID] = user

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save user data: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*types.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user with email %s not found", email)
}

// GetAll retrieves all users
func (r *userRepository) GetAll(ctx context.Context) ([]*types.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*types.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *types.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user exists
	existingUser, exists := r.users[user.ID]
	if !exists {
		return fmt.Errorf("user with ID %s not found", user.ID)
	}

	// Check if email is being changed and if new email already exists
	if user.Email != existingUser.Email {
		for _, u := range r.users {
			if u.ID != user.ID && u.Email == user.Email {
				return fmt.Errorf("user with email %s already exists", user.Email)
			}
		}
	}

	// Update timestamp
	user.UpdatedAt = time.Now()
	user.CreatedAt = existingUser.CreatedAt

	// Preserve IssueIDs if not provided
	if user.IssueIDs == nil {
		user.IssueIDs = existingUser.IssueIDs
	}

	r.users[user.ID] = user

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save user data: %w", err)
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user with ID %s not found", id)
	}

	delete(r.users, id)

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save user data: %w", err)
	}

	return nil
}

// loadData loads user data from JSON file
func (r *userRepository) loadData() error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		// File doesn't exist, create empty data
		return r.saveData()
	}

	// Read file
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read user data file: %w", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, &r.users); err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return nil
}

// saveData saves user data to JSON file
func (r *userRepository) saveData() error {
	// Marshal to JSON
	data, err := json.MarshalIndent(r.users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write user data file: %w", err)
	}

	return nil
}
