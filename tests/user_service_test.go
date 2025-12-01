package tests

import (
	"context"
	"testing"

	"jira-clone-be/internal/service"
	"jira-clone-be/pkg/logger"
	"jira-clone-be/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger := logger.New("info")
	userService := service.NewUserService(mockRepo, logger)

	// Test data
	req := &types.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, assert.AnError)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*types.User")).Return(nil)

	// Execute
	user, err := userService.CreateUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger := logger.New("info")
	userService := service.NewUserService(mockRepo, logger)

	// Test data
	userID := "test-user-id"
	expectedUser := &types.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	mockRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	// Execute
	user, err := userService.GetUser(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger := logger.New("info")
	userService := service.NewUserService(mockRepo, logger)

	// Test data
	userID := "test-user-id"
	req := &types.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	existingUser := &types.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*types.User")).Return(nil)

	// Execute
	user, err := userService.UpdateUser(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, userID, user.ID)

	mockRepo.AssertExpectations(t)
}
