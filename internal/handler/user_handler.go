package handler

import (
	"encoding/json"
	"net/http"

	"jira-clone-be/types"
	"jira-clone-be/internal/service"
	"jira-clone-be/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles user requests
type UserHandler struct {
	userService *service.UserService
	logger      *logger.Logger
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
		validator:   validator.New(),
	}
}

// GetUsers retrieves all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetUsers(r.Context())
	if err != nil {
		h.logger.Error("Failed to get users", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, users)
}

// GetUser retrieves a specific user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, user)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var req types.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode update user request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Update user request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to update user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}
