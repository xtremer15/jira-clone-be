package types

import "time"

// Issue represents an issue/task in the system with all its associated metadata and workflow information.
// It contains the core issue information including priority, status, assignment, and project association.
type Issue struct {
	ID          string        `json:"id" validate:"required"`                                             // Unique identifier for the issue
	Priority    IssuePriority `json:"priority" validate:"required,oneof=LOWEST LOW MEDIUM HIGH HIGHEST"`  // Issue priority level
	Status      IssueStatus   `json:"status" validate:"required,oneof=BACKLOG SELECTED IN_PROGRESS DONE"` // Current workflow status
	Title       string        `json:"title" validate:"required,min=1,max=200"`                            // Issue title (1-200 characters)
	Description string        `json:"description" validate:"max=1000"`                                    // Issue description (maximum 1000 characters)
	Type        IssueType     `json:"type" validate:"required,oneof=BUG STORY TASK"`                      // Type of issue (bug, story, or task)
	Icon        string        `json:"icon"`                                                               // Icon identifier for the issue type
	Assignee    string        `json:"assignee"`                                                           // User ID of the person assigned to this issue
	Reporter    string        `json:"reporter"`                                                           // User ID of the person who reported this issue
	IconURL     string        `json:"iconUrl"`                                                            // URL to the issue icon image
	DisplayID   string        `json:"displayId" validate:"required"`                                      // Human-readable issue identifier (e.g., "PROJ-123")
	ProjectID   string        `json:"projectId" validate:"required"`                                      // ID of the project this issue belongs to
	CreatedAt   time.Time     `json:"createdAt"`                                                          // Timestamp when the issue was created
	UpdatedAt   time.Time     `json:"updatedAt"`                                                          // Timestamp when the issue was last updated
}

// CreateIssueRequest represents the request payload for creating a new issue.
// It contains all the essential information needed to create an issue in the system.
type CreateIssueRequest struct {
	Priority    IssuePriority `json:"priority" validate:"required,oneof=LOWEST LOW MEDIUM HIGH HIGHEST"`  // Issue priority level
	Status      IssueStatus   `json:"status" validate:"required,oneof=BACKLOG SELECTED IN_PROGRESS DONE"` // Initial workflow status
	Title       string        `json:"title" validate:"required,min=1,max=200"`                            // Issue title (1-200 characters)
	Description string        `json:"description" validate:"max=1000"`                                    // Issue description (maximum 1000 characters)
	Type        IssueType     `json:"type" validate:"required,oneof=BUG STORY TASK"`                      // Type of issue (bug, story, or task)
	Assignee    string        `json:"assignee"`                                                           // User ID of the person assigned to this issue
	Reporter    string        `json:"reporter"`                                                           // User ID of the person who reported this issue
	ProjectID   string        `json:"projectId" validate:"required"`                                      // ID of the project this issue belongs to
}

// UpdateIssueRequest represents the request payload for updating an existing issue.
// All fields are optional, allowing partial updates of issue information.
type UpdateIssueRequest struct {
	Priority    IssuePriority `json:"priority" validate:"omitempty,oneof=LOWEST LOW MEDIUM HIGH HIGHEST"`  // Updated issue priority level (optional)
	Status      IssueStatus   `json:"status" validate:"omitempty,oneof=BACKLOG SELECTED IN_PROGRESS DONE"` // Updated workflow status (optional)
	Title       string        `json:"title" validate:"omitempty,min=1,max=200"`                            // Updated issue title (1-200 characters, optional)
	Description string        `json:"description" validate:"omitempty,max=1000"`                           // Updated issue description (maximum 1000 characters, optional)
	Type        IssueType     `json:"type" validate:"omitempty,oneof=BUG STORY TASK"`                      // Updated issue type (optional)
	Assignee    string        `json:"assignee"`                                                            // Updated assignee user ID (optional)
	Reporter    string        `json:"reporter"`                                                            // Updated reporter user ID (optional)
}
