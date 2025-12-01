package service

import (
	"context"
	"fmt"
	"time"

	"jira-clone-be/types"
	"jira-clone-be/internal/repository"
	"jira-clone-be/pkg/logger"

	"github.com/google/uuid"
)

// ProjectService handles project business logic
type ProjectService struct {
	projectRepo repository.ProjectRepository
	logger      *logger.Logger
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo repository.ProjectRepository, logger *logger.Logger) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		logger:      logger,
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error) {
	// Create project
	project := &types.Project{
		ID:          uuid.New().String(),
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Category:    req.Category,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Issues:      []types.Issue{},
		Users:       []types.User{},
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		s.logger.Error("Failed to create project", map[string]interface{}{
			"error": err.Error(),
			"name":  req.Name,
		})
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	s.logger.Info("Project created successfully", map[string]interface{}{
		"projectID": project.ID,
		"name":      project.Name,
	})

	return project, nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id string) (*types.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get project", map[string]interface{}{
			"error":     err.Error(),
			"projectID": id,
		})
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// GetProjects retrieves all projects
func (s *ProjectService) GetProjects(ctx context.Context) ([]*types.Project, error) {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get projects", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	return projects, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, id string, req *types.UpdateProjectRequest) (*types.Project, error) {
	// Get existing project
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get project for update", map[string]interface{}{
			"error":     err.Error(),
			"projectID": id,
		})
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Update fields
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.URL != "" {
		project.URL = req.URL
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Category != "" {
		project.Category = req.Category
	}
	project.UpdatedAt = time.Now()

	// Save updated project
	if err := s.projectRepo.Update(ctx, project); err != nil {
		s.logger.Error("Failed to update project", map[string]interface{}{
			"error":     err.Error(),
			"projectID": id,
		})
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	s.logger.Info("Project updated successfully", map[string]interface{}{
		"projectID": project.ID,
		"name":      project.Name,
	})

	return project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get project for deletion", map[string]interface{}{
			"error":     err.Error(),
			"projectID": id,
		})
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Delete project
	if err := s.projectRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete project", map[string]interface{}{
			"error":     err.Error(),
			"projectID": id,
		})
		return fmt.Errorf("failed to delete project: %w", err)
	}

	s.logger.Info("Project deleted successfully", map[string]interface{}{
		"projectID": id,
	})

	return nil
}

// GetProjectBoard retrieves the board view for a project
func (s *ProjectService) GetProjectBoard(ctx context.Context, projectID string) (*types.Board, error) {
	// Get project
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		s.logger.Error("Failed to get project for board", map[string]interface{}{
			"error":     err.Error(),
			"projectID": projectID,
		})
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Create board with lanes
	board := &types.Board{
		ProjectID: projectID,
		Lanes: []types.Lane{
			{
				ID:     types.StatusBacklog,
				Title:  "Backlog",
				Issues: []types.Issue{},
			},
			{
				ID:     types.StatusSelected,
				Title:  "Selected",
				Issues: []types.Issue{},
			},
			{
				ID:     types.StatusInProgress,
				Title:  "In Progress",
				Issues: []types.Issue{},
			},
			{
				ID:     types.StatusDone,
				Title:  "Done",
				Issues: []types.Issue{},
			},
		},
	}

	// Group issues by status
	for _, issue := range project.Issues {
		for i, lane := range board.Lanes {
			if lane.ID == issue.Status {
				board.Lanes[i].Issues = append(board.Lanes[i].Issues, issue)
				break
			}
		}
	}

	return board, nil
}
