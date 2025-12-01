package service

import (
	"context"
	"fmt"
	"time"

	"jira-clone-be/internal/repository"
	"jira-clone-be/pkg/logger"
	"jira-clone-be/types"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	logger    *logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(ctx context.Context, req *types.LoginRequest) (*types.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Failed to get user for login", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		return nil, fmt.Errorf("invalid credentials")
	}

	// In a real implementation, you'd validate the password against the stored hash
	// For now, we'll just check if the user exists
	// TODO: Implement proper password validation

	// Generate JWT tokens
	accessToken, refreshToken, expiresIn, err := s.generateTokens(user.ID, user.Email)
	if err != nil {
		s.logger.Error("Failed to generate tokens", map[string]interface{}{
			"error":  err.Error(),
			"userID": user.ID,
		})
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info("User logged in successfully", map[string]interface{}{
		"userID": user.ID,
		"email":  user.Email,
	})

	return &types.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// Register creates a new user and returns JWT tokens
func (s *AuthService) Register(ctx context.Context, req *types.RegisterRequest) (*types.AuthResponse, error) {
	fmt.Println("Registering user", req)
	fmt.Println("Registering user context", ctx)
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
		ID:        fmt.Sprintf("user_%d", time.Now().Unix()),
		Name:      req.Name,
		Email:     req.Email,
		AvatarURL: fmt.Sprintf("/assets/avatars/%s.png", req.Name),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IssueIDs:  []string{},
	}

	// Store user (in a real implementation, you'd also store the hashed password)
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user during registration", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT tokens
	accessToken, refreshToken, expiresIn, err := s.generateTokens(user.ID, user.Email)
	if err != nil {
		s.logger.Error("Failed to generate tokens", map[string]interface{}{
			"error":  err.Error(),
			"userID": user.ID,
		})
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info("User registered successfully", map[string]interface{}{
		"userID": user.ID,
		"email":  user.Email,
	})

	return &types.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*types.AuthResponse, error) {
	// Parse and validate refresh token
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		s.logger.Error("Failed to parse refresh token", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims["userID"].(string))
	if err != nil {
		s.logger.Error("Failed to get user for token refresh", map[string]interface{}{
			"error":  err.Error(),
			"userID": claims["userID"],
		})
		return nil, fmt.Errorf("user not found")
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresIn, err := s.generateTokens(user.ID, user.Email)
	if err != nil {
		s.logger.Error("Failed to generate new tokens", map[string]interface{}{
			"error":  err.Error(),
			"userID": user.ID,
		})
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info("Tokens refreshed successfully", map[string]interface{}{
		"userID": user.ID,
		"email":  user.Email,
	})

	return &types.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return userID, nil
}

// generateTokens generates access and refresh tokens
func (s *AuthService) generateTokens(userID, email string) (string, string, int64, error) {
	// Access token (expires in 15 minutes)
	accessClaims := jwt.MapClaims{
		"userID": userID,
		"email":  email,
		"exp":    time.Now().Add(60 * time.Minute).Unix(),
		"type":   "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token (expires in 7 days)
	refreshClaims := jwt.MapClaims{
		"userID": userID,
		"email":  email,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type":   "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	expiresIn := int64(60 * 60) // 15 minutes in seconds

	return accessTokenString, refreshTokenString, expiresIn, nil
}

// parseToken parses and validates a JWT token
func (s *AuthService) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
