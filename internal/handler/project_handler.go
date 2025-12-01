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

// ProjectHandler handles project requests
type ProjectHandler struct {
	projectService *service.ProjectService
	logger         *logger.Logger
	validator      *validator.Validate
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService *service.ProjectService, logger *logger.Logger) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		logger:         logger,
		validator:      validator.New(),
	}
}

// GetProjects retrieves all projects
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetProjects(r.Context())
	if err != nil {
		h.logger.Error("Failed to get projects", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to get projects", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, projects)
}

// CreateProject creates a new project
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req types.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode create project request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Create project request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Create project
	project, err := h.projectService.CreateProject(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create project", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, project)
}

// GetProject retrieves a specific project
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	project, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		h.logger.Error("Failed to get project", map[string]interface{}{
			"error":     err.Error(),
			"projectID": projectID,
		})
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, project)
}

// UpdateProject updates a project
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	var req types.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode update project request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Update project request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Update project
	project, err := h.projectService.UpdateProject(r.Context(), projectID, &req)
	if err != nil {
		h.logger.Error("Failed to update project", map[string]interface{}{
			"error":     err.Error(),
			"projectID": projectID,
		})
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, project)
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	err := h.projectService.DeleteProject(r.Context(), projectID)
	if err != nil {
		h.logger.Error("Failed to delete project", map[string]interface{}{
			"error":     err.Error(),
			"projectID": projectID,
		})
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
