package handler

import (
	"encoding/json"
	"net/http"

	"jira-clone-be/types"
	"jira-clone-be/internal/service"
	"jira-clone-be/pkg/logger"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
	validator   *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		validator:   validator.New(),
	}
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode login request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Login request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Authenticate user
	response, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		h.logger.Error("Login failed", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Return response
	render.JSON(w, r, response)
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode register request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Register request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Register user
	response, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		h.logger.Error("Registration failed", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}

	// Return response
	render.JSON(w, r, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode refresh token request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Refresh token request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Refresh token
	response, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("Token refresh failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Return response
	render.JSON(w, r, response)
}
