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

// IssueService handles issue business logic
type IssueService struct {
	issueRepo repository.IssueRepository
	logger    *logger.Logger
}

// NewIssueService creates a new issue service
func NewIssueService(issueRepo repository.IssueRepository, logger *logger.Logger) *IssueService {
	return &IssueService{
		issueRepo: issueRepo,
		logger:    logger,
	}
}

// CreateIssue creates a new issue
func (s *IssueService) CreateIssue(ctx context.Context, req *types.CreateIssueRequest) (*types.Issue, error) {
	// Create issue
	issue := &types.Issue{
		ID:          uuid.New().String(),
		Priority:    req.Priority,
		Status:      req.Status,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Icon:        s.getIconForType(req.Type),
		Assignee:    req.Assignee,
		Reporter:    req.Reporter,
		IconURL:     s.getIconURLForType(req.Type),
		DisplayID:   s.generateDisplayID(),
		ProjectID:   req.ProjectID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.issueRepo.Create(ctx, issue); err != nil {
		s.logger.Error("Failed to create issue", map[string]interface{}{
			"error": err.Error(),
			"title": req.Title,
		})
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	s.logger.Info("Issue created successfully", map[string]interface{}{
		"issueID":   issue.ID,
		"displayID": issue.DisplayID,
		"title":     issue.Title,
	})

	return issue, nil
}

// GetIssue retrieves an issue by ID
func (s *IssueService) GetIssue(ctx context.Context, id string) (*types.Issue, error) {
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get issue", map[string]interface{}{
			"error":   err.Error(),
			"issueID": id,
		})
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return issue, nil
}

// GetIssues retrieves all issues
func (s *IssueService) GetIssues(ctx context.Context) ([]*types.Issue, error) {
	issues, err := s.issueRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get issues", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	return issues, nil
}

// GetIssuesByProject retrieves all issues for a project
func (s *IssueService) GetIssuesByProject(ctx context.Context, projectID string) ([]*types.Issue, error) {
	issues, err := s.issueRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		s.logger.Error("Failed to get issues by project", map[string]interface{}{
			"error":     err.Error(),
			"projectID": projectID,
		})
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	return issues, nil
}

// GetIssuesByStatus retrieves issues by status
func (s *IssueService) GetIssuesByStatus(ctx context.Context, status types.IssueStatus) ([]*types.Issue, error) {
	issues, err := s.issueRepo.GetByStatus(ctx, status)
	if err != nil {
		s.logger.Error("Failed to get issues by status", map[string]interface{}{
			"error":  err.Error(),
			"status": status,
		})
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	return issues, nil
}

// GetIssuesByAssignee retrieves issues by assignee
func (s *IssueService) GetIssuesByAssignee(ctx context.Context, assignee string) ([]*types.Issue, error) {
	issues, err := s.issueRepo.GetByAssignee(ctx, assignee)
	if err != nil {
		s.logger.Error("Failed to get issues by assignee", map[string]interface{}{
			"error":    err.Error(),
			"assignee": assignee,
		})
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	return issues, nil
}

// UpdateIssue updates an existing issue
func (s *IssueService) UpdateIssue(ctx context.Context, id string, req *types.UpdateIssueRequest) (*types.Issue, error) {
	// Get existing issue
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get issue for update", map[string]interface{}{
			"error":   err.Error(),
			"issueID": id,
		})
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	// Update fields
	if req.Priority != "" {
		issue.Priority = req.Priority
	}
	if req.Status != "" {
		issue.Status = req.Status
	}
	if req.Title != "" {
		issue.Title = req.Title
	}
	if req.Description != "" {
		issue.Description = req.Description
	}
	if req.Type != "" {
		issue.Type = req.Type
		issue.Icon = s.getIconForType(req.Type)
		issue.IconURL = s.getIconURLForType(req.Type)
	}
	if req.Assignee != "" {
		issue.Assignee = req.Assignee
	}
	if req.Reporter != "" {
		issue.Reporter = req.Reporter
	}
	issue.UpdatedAt = time.Now()

	// Save updated issue
	if err := s.issueRepo.Update(ctx, issue); err != nil {
		s.logger.Error("Failed to update issue", map[string]interface{}{
			"error":   err.Error(),
			"issueID": id,
		})
		return nil, fmt.Errorf("failed to update issue: %w", err)
	}

	s.logger.Info("Issue updated successfully", map[string]interface{}{
		"issueID":   issue.ID,
		"displayID": issue.DisplayID,
		"title":     issue.Title,
	})

	return issue, nil
}

// DeleteIssue deletes an issue
func (s *IssueService) DeleteIssue(ctx context.Context, id string) error {
	// Check if issue exists
	_, err := s.issueRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get issue for deletion", map[string]interface{}{
			"error":   err.Error(),
			"issueID": id,
		})
		return fmt.Errorf("failed to get issue: %w", err)
	}

	// Delete issue
	if err := s.issueRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete issue", map[string]interface{}{
			"error":   err.Error(),
			"issueID": id,
		})
		return fmt.Errorf("failed to delete issue: %w", err)
	}

	s.logger.Info("Issue deleted successfully", map[string]interface{}{
		"issueID": id,
	})

	return nil
}

// getIconForType returns the icon for an issue type
func (s *IssueService) getIconForType(issueType types.IssueType) string {
	switch issueType {
	case types.TypeBug:
		return "bug"
	case types.TypeStory:
		return "story"
	case types.TypeTask:
		return "task"
	default:
		return "task"
	}
}

// getIconURLForType returns the icon URL for an issue type
func (s *IssueService) getIconURLForType(issueType types.IssueType) string {
	switch issueType {
	case types.TypeBug:
		return "/assets/icons/bug.png"
	case types.TypeStory:
		return "/assets/icons/story.png"
	case types.TypeTask:
		return "/assets/icons/task.png"
	default:
		return "/assets/icons/task.png"
	}
}

// generateDisplayID generates a unique display ID for an issue
func (s *IssueService) generateDisplayID() string {
	// In a real implementation, you'd generate a proper display ID
	// For now, we'll use a simple counter
	return fmt.Sprintf("ISS-%d", time.Now().Unix())
}
