package repository

import (
	"context"
	"jira-clone-be/types"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *types.User) error
	GetByID(ctx context.Context, id string) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	GetAll(ctx context.Context) ([]*types.User, error)
	Update(ctx context.Context, user *types.User) error
	Delete(ctx context.Context, id string) error
}

// ProjectRepository defines the interface for project data operations
type ProjectRepository interface {
	Create(ctx context.Context, project *types.Project) error
	GetByID(ctx context.Context, id string) (*types.Project, error)
	GetAll(ctx context.Context) ([]*types.Project, error)
	Update(ctx context.Context, project *types.Project) error
	Delete(ctx context.Context, id string) error
}

// IssueRepository defines the interface for issue data operations
type IssueRepository interface {
	Create(ctx context.Context, issue *types.Issue) error
	GetByID(ctx context.Context, id string) (*types.Issue, error)
	GetByProjectID(ctx context.Context, projectID string) ([]*types.Issue, error)
	GetAll(ctx context.Context) ([]*types.Issue, error)
	Update(ctx context.Context, issue *types.Issue) error
	Delete(ctx context.Context, id string) error
	GetByStatus(ctx context.Context, status types.IssueStatus) ([]*types.Issue, error)
	GetByAssignee(ctx context.Context, assignee string) ([]*types.Issue, error)
}
