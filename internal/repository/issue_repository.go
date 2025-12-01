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

// issueRepository implements IssueRepository interface using JSON file storage
type issueRepository struct {
	filePath string
	issues   map[string]*types.Issue
	mutex    sync.RWMutex
	logger   *logger.Logger
}

// NewIssueRepository creates a new issue repository instance
func NewIssueRepository() IssueRepository {
	repo := &issueRepository{
		filePath: "data/issues.json",
		issues:   make(map[string]*types.Issue),
		logger:   logger.New("info"),
	}

	// Load existing data
	if err := repo.loadData(); err != nil {
		repo.logger.Warn("Failed to load issue data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return repo
}

// Create creates a new issue
func (r *issueRepository) Create(ctx context.Context, issue *types.Issue) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate ID if not provided
	if issue.ID == "" {
		issue.ID = uuid.New().String()
	}

	// Generate display ID if not provided
	if issue.DisplayID == "" {
		issue.DisplayID = fmt.Sprintf("ISS-%d", len(r.issues)+1)
	}

	// Set timestamps
	now := time.Now()
	issue.CreatedAt = now
	issue.UpdatedAt = now

	r.issues[issue.ID] = issue

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save issue data: %w", err)
	}

	return nil
}

// GetByID retrieves an issue by ID
func (r *issueRepository) GetByID(ctx context.Context, id string) (*types.Issue, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	issue, exists := r.issues[id]
	if !exists {
		return nil, fmt.Errorf("issue with ID %s not found", id)
	}

	return issue, nil
}

// GetByProjectID retrieves all issues for a project
func (r *issueRepository) GetByProjectID(ctx context.Context, projectID string) ([]*types.Issue, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var issues []*types.Issue
	for _, issue := range r.issues {
		if issue.ProjectID == projectID {
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// GetAll retrieves all issues
func (r *issueRepository) GetAll(ctx context.Context) ([]*types.Issue, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	issues := make([]*types.Issue, 0, len(r.issues))
	for _, issue := range r.issues {
		issues = append(issues, issue)
	}

	return issues, nil
}

// Update updates an existing issue
func (r *issueRepository) Update(ctx context.Context, issue *types.Issue) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if issue exists
	existingIssue, exists := r.issues[issue.ID]
	if !exists {
		return fmt.Errorf("issue with ID %s not found", issue.ID)
	}

	// Update timestamp
	issue.UpdatedAt = time.Now()
	issue.CreatedAt = existingIssue.CreatedAt

	// Preserve DisplayID if not provided
	if issue.DisplayID == "" {
		issue.DisplayID = existingIssue.DisplayID
	}

	r.issues[issue.ID] = issue

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save issue data: %w", err)
	}

	return nil
}

// Delete deletes an issue by ID
func (r *issueRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.issues[id]; !exists {
		return fmt.Errorf("issue with ID %s not found", id)
	}

	delete(r.issues, id)

	// Save to file
	if err := r.saveData(); err != nil {
		return fmt.Errorf("failed to save issue data: %w", err)
	}

	return nil
}

// GetByStatus retrieves issues by status
func (r *issueRepository) GetByStatus(ctx context.Context, status types.IssueStatus) ([]*types.Issue, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var issues []*types.Issue
	for _, issue := range r.issues {
		if issue.Status == status {
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// GetByAssignee retrieves issues by assignee
func (r *issueRepository) GetByAssignee(ctx context.Context, assignee string) ([]*types.Issue, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var issues []*types.Issue
	for _, issue := range r.issues {
		if issue.Assignee == assignee {
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// loadData loads issue data from JSON file
func (r *issueRepository) loadData() error {
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
		return fmt.Errorf("failed to read issue data file: %w", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, &r.issues); err != nil {
		return fmt.Errorf("failed to unmarshal issue data: %w", err)
	}

	return nil
}

// saveData saves issue data to JSON file
func (r *issueRepository) saveData() error {
	// Marshal to JSON
	data, err := json.MarshalIndent(r.issues, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal issue data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write issue data file: %w", err)
	}

	return nil
}
