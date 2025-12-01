// Package types provides domain models, DTOs, and enums for the Jira Clone Backend application.
// It contains all the core data structures used throughout the application.
package types

// IssuePriority represents the priority level of an issue in the project management system.
// Higher priority issues should be addressed before lower priority ones.
type IssuePriority string

const (
	PriorityLowest  IssuePriority = "LOWEST"  // Lowest priority - can be addressed last
	PriorityLow     IssuePriority = "LOW"     // Low priority - not urgent
	PriorityMedium  IssuePriority = "MEDIUM"  // Medium priority - standard importance
	PriorityHigh    IssuePriority = "HIGH"    // High priority - important and urgent
	PriorityHighest IssuePriority = "HIGHEST" // Highest priority - critical and must be addressed first
)

// IssueStatus represents the current status of an issue in the workflow.
// Issues progress through these statuses from creation to completion.
type IssueStatus string

const (
	StatusBacklog    IssueStatus = "BACKLOG"     // Issue is in the backlog, not yet selected for work
	StatusSelected   IssueStatus = "SELECTED"    // Issue has been selected for the current sprint/iteration
	StatusInProgress IssueStatus = "IN_PROGRESS" // Issue is currently being worked on
	StatusDone       IssueStatus = "DONE"        // Issue has been completed
)

// IssueType represents the type of an issue, which determines its purpose and workflow.
// Different types may have different fields, workflows, or validation rules.
type IssueType string

const (
	TypeBug   IssueType = "BUG"   // Bug report - something that is broken and needs to be fixed
	TypeStory IssueType = "STORY" // User story - a feature or functionality requested by users
	TypeTask  IssueType = "TASK"  // Task - a specific piece of work to be completed
)

// ProjectCategory represents the category of a project, used for organization and filtering.
// Projects can be grouped by category for better management and reporting.
type ProjectCategory string

const (
	CategorySoftware  ProjectCategory = "Software"  // Software development projects
	CategoryMarketing ProjectCategory = "Marketing" // Marketing and promotional projects
	CategoryBusiness  ProjectCategory = "Business"  // Business process and operational projects
)

// AvatarSize represents the size of a user avatar image.
// Different sizes are used in different parts of the UI for optimal display.
type AvatarSize string

const (
	AvatarSizeLarge   AvatarSize = "large"   // Large avatar for profile pages and detailed views
	AvatarSizeSmall   AvatarSize = "small"   // Small avatar for lists and compact views
	AvatarSizeDefault AvatarSize = "default" // Default size for general use
)

// AvatarShape represents the shape of a user avatar image.
// Different shapes provide different visual styles for the user interface.
type AvatarShape string

const (
	AvatarShapeCircle AvatarShape = "circle" // Circular avatar for a modern, friendly look
	AvatarShapeSquare AvatarShape = "square" // Square avatar for a more formal appearance
)
