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

// IssueHandler handles issue requests
type IssueHandler struct {
	issueService *service.IssueService
	logger       *logger.Logger
	validator    *validator.Validate
}

// NewIssueHandler creates a new issue handler
func NewIssueHandler(issueService *service.IssueService, logger *logger.Logger) *IssueHandler {
	return &IssueHandler{
		issueService: issueService,
		logger:       logger,
		validator:    validator.New(),
	}
}

// GetIssues retrieves all issues
func (h *IssueHandler) GetIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := h.issueService.GetIssues(r.Context())
	if err != nil {
		h.logger.Error("Failed to get issues", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to get issues", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, issues)
}

// CreateIssue creates a new issue
func (h *IssueHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var req types.CreateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode create issue request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Create issue request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Create issue
	issue, err := h.issueService.CreateIssue(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create issue", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to create issue", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, issue)
}

// GetIssue retrieves a specific issue
func (h *IssueHandler) GetIssue(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")
	if issueID == "" {
		http.Error(w, "Issue ID is required", http.StatusBadRequest)
		return
	}

	issue, err := h.issueService.GetIssue(r.Context(), issueID)
	if err != nil {
		h.logger.Error("Failed to get issue", map[string]interface{}{
			"error":   err.Error(),
			"issueID": issueID,
		})
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, issue)
}

// UpdateIssue updates an issue
func (h *IssueHandler) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")
	if issueID == "" {
		http.Error(w, "Issue ID is required", http.StatusBadRequest)
		return
	}

	var req types.UpdateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode update issue request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Update issue request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Update issue
	issue, err := h.issueService.UpdateIssue(r.Context(), issueID, &req)
	if err != nil {
		h.logger.Error("Failed to update issue", map[string]interface{}{
			"error":   err.Error(),
			"issueID": issueID,
		})
		http.Error(w, "Failed to update issue", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, issue)
}

// DeleteIssue deletes an issue
func (h *IssueHandler) DeleteIssue(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")
	if issueID == "" {
		http.Error(w, "Issue ID is required", http.StatusBadRequest)
		return
	}

	err := h.issueService.DeleteIssue(r.Context(), issueID)
	if err != nil {
		h.logger.Error("Failed to delete issue", map[string]interface{}{
			"error":   err.Error(),
			"issueID": issueID,
		})
		http.Error(w, "Failed to delete issue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
