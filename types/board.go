package types

// Board represents a board view of a project, organized into lanes/columns.
// It provides a Kanban-style view of project issues grouped by their status.
type Board struct {
	ProjectID string `json:"projectId" validate:"required"` // ID of the project this board belongs to
	Lanes     []Lane `json:"lanes"`                         // List of lanes/columns in the board
}

// Lane represents a column/lane in the board view.
// Each lane corresponds to a specific issue status and contains the issues in that status.
type Lane struct {
	ID     IssueStatus `json:"id" validate:"required"`    // Status identifier for this lane
	Title  string      `json:"title" validate:"required"` // Display title for this lane
	Issues []Issue     `json:"issues"`                    // List of issues in this lane/status
}

// GetBoardRequest represents the request payload for retrieving a project board.
// It specifies which project's board view to retrieve.
type GetBoardRequest struct {
	ProjectID string `json:"projectId" validate:"required"` // ID of the project to get the board for
}
