package types

import "time"

// Project represents a project in the system with all its associated data and metadata.
// It contains project information, associated issues, and team members.
type Project struct {
	ID          string          `json:"id" validate:"required"`                                         // Unique identifier for the project
	Name        string          `json:"name" validate:"required,min=1,max=100"`                         // Project name (1-100 characters)
	URL         string          `json:"url" validate:"required,url"`                                    // Project URL (must be valid URL format)
	Description string          `json:"description" validate:"max=500"`                                 // Project description (maximum 500 characters)
	Category    ProjectCategory `json:"category" validate:"required,oneof=Software Marketing Business"` // Project category
	CreatedAt   time.Time       `json:"createdAt"`                                                      // Timestamp when the project was created
	UpdatedAt   time.Time       `json:"updatedAt"`                                                      // Timestamp when the project was last updated
	Issues      []Issue         `json:"issues"`                                                         // List of issues associated with this project
	Users       []User          `json:"users"`                                                          // List of users who have access to this project
}

// CreateProjectRequest represents the request payload for creating a new project.
// It contains all the essential information needed to create a project.
type CreateProjectRequest struct {
	Name        string          `json:"name" validate:"required,min=1,max=100"`                         // Project name (1-100 characters)
	URL         string          `json:"url" validate:"required,url"`                                    // Project URL (must be valid URL format)
	Description string          `json:"description" validate:"max=500"`                                 // Project description (maximum 500 characters)
	Category    ProjectCategory `json:"category" validate:"required,oneof=Software Marketing Business"` // Project category
}

// UpdateProjectRequest represents the request payload for updating an existing project.
// All fields are optional, allowing partial updates of project information.
type UpdateProjectRequest struct {
	Name        string          `json:"name" validate:"omitempty,min=1,max=100"`                         // Updated project name (1-100 characters, optional)
	URL         string          `json:"url" validate:"omitempty,url"`                                    // Updated project URL (valid URL format, optional)
	Description string          `json:"description" validate:"omitempty,max=500"`                        // Updated project description (maximum 500 characters, optional)
	Category    ProjectCategory `json:"category" validate:"omitempty,oneof=Software Marketing Business"` // Updated project category (optional)
}
