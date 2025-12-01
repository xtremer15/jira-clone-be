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

// projectRepository implements ProjectRepository interface using JSON file storage
type projectRepository struct {
	filePath string
	projects map[string]*types.Project
	mutex    sync.RWMutex
	logger   *logger.Logger
}

// NewProjectRepository creates a new project repository instance
func NewProjectRepository() ProjectRepository {
	repo := &projectRepository{
		filePath: "data/projects.json",
		projects: make(map[string]*types.Project),
		logger:   logger.New("info"),
	}

	// Load existing data
	if err := repo.loadData(); err != nil {
		repo.logger.Warn("Failed to load project data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return repo
}

// Create creates a new project
func (r *projectRepository) Create(ctx context.Context, project *types.Project) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate ID if not provided
	if project.ID == "" {
		project.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	// Initialize empty slices
	if project.Issues == nil {
		project.Issues = []types.Issue{}
	}
	if project.Users == nil {
		project.Users = []types.User{}
	}

	r.projects[project.ID] = project

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save project data: %w", err)
	}

	return nil
}

// GetByID retrieves a project by ID
func (r *projectRepository) GetByID(ctx context.Context, id string) (*types.Project, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	project, exists := r.projects[id]
	if !exists {
		return nil, fmt.Errorf("project with ID %s not found", id)
	}

	return project, nil
}

// GetAll retrieves all projects
func (r *projectRepository) GetAll(ctx context.Context) ([]*types.Project, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	projects := make([]*types.Project, 0, len(r.projects))
	for _, project := range r.projects {
		projects = append(projects, project)
	}

	return projects, nil
}

// Update updates an existing project
func (r *projectRepository) Update(ctx context.Context, project *types.Project) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if project exists
	existingProject, exists := r.projects[project.ID]
	if !exists {
		return fmt.Errorf("project with ID %s not found", project.ID)
	}

	// Update timestamp
	project.UpdatedAt = time.Now()
	project.CreatedAt = existingProject.CreatedAt

	// Preserve Issues and Users if not provided
	if project.Issues == nil {
		project.Issues = existingProject.Issues
	}
	if project.Users == nil {
		project.Users = existingProject.Users
	}

	r.projects[project.ID] = project

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save project data: %w", err)
	}

	return nil
}

// Delete deletes a project by ID
func (r *projectRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.projects[id]; !exists {
		return fmt.Errorf("project with ID %s not found", id)
	}

	delete(r.projects, id)

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save project data: %w", err)
	}

	return nil
}

// loadData loads project data from JSON file
func (r *projectRepository) loadData() error {
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
		return fmt.Errorf("failed to read project data file: %w", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, &r.projects); err != nil {
		return fmt.Errorf("failed to unmarshal project data: %w", err)
	}

	return nil
}

// saveData saves project data to JSON file
func (r *projectRepository) saveData() error {
	// Marshal to JSON
	data, err := json.MarshalIndent(r.projects, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project data file: %w", err)
	}

	return nil
}
