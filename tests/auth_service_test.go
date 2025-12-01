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

func TestAuthService_Login(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger := logger.New("info")
	authService := service.NewAuthService(mockRepo, "test-secret", logger)

	// Test data
	req := &types.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	expectedUser := &types.User{
		ID:    "test-user-id",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(expectedUser, nil)

	// Execute
	response, err := authService.Login(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.User)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Greater(t, response.ExpiresIn, int64(0))
	assert.Equal(t, expectedUser.ID, response.User.ID)
	assert.Equal(t, expectedUser.Email, response.User.Email)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger := logger.New("info")
	authService := service.NewAuthService(mockRepo, "test-secret", logger)

	// Test data
	req := &types.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, assert.AnError)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*types.User")).Return(nil)

	// Execute
	response, err := authService.Register(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.User)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Greater(t, response.ExpiresIn, int64(0))
	assert.Equal(t, req.Name, response.User.Name)
	assert.Equal(t, req.Email, response.User.Email)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger := logger.New("info")
	authService := service.NewAuthService(mockRepo, "test-secret", logger)

	// Test data - using an invalid token to test error handling
	accessToken := "invalid-token"

	// Execute
	validatedUserID, err := authService.ValidateToken(accessToken)

	// Assert - should fail with invalid token
	assert.Error(t, err)
	assert.Empty(t, validatedUserID)
}

func TestAuthService_RefreshToken(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	logger := logger.New("info")
	authService := service.NewAuthService(mockRepo, "test-secret", logger)

	// Test data - using an invalid refresh token to test error handling
	refreshToken := "invalid-refresh-token"

	// Execute
	response, err := authService.RefreshToken(context.Background(), refreshToken)

	// Assert - should fail with invalid token
	assert.Error(t, err)
	assert.Nil(t, response)
}
