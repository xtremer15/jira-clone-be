package types

import "time"

// User represents a user in the system with all their profile information and associated data.
// It contains the core user information including identification, contact details, and metadata.
type User struct {
	ID        string    `json:"id" validate:"required"`                 // Unique identifier for the user
	Name      string    `json:"name" validate:"required,min=1,max=100"` // User's full name
	Email     string    `json:"email" validate:"required,email"`        // User's email address (unique)
	AvatarURL string    `json:"avatarUrl"`                              // URL to user's profile picture
	CreatedAt time.Time `json:"createdAt"`                              // Timestamp when the user was created
	UpdatedAt time.Time `json:"updatedAt"`                              // Timestamp when the user was last updated
	IssueIDs  []string  `json:"issueIds"`                               // List of issue IDs assigned to this user
}

// CreateUserRequest represents the request payload for creating a new user.
// It contains the essential information needed to register a new user in the system.
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"` // User's full name (1-100 characters)
	Email    string `json:"email" validate:"required,email"`        // User's email address (must be valid email format)
	Password string `json:"password" validate:"required,min=6"`     // User's password (minimum 6 characters)
}

// UpdateUserRequest represents the request payload for updating an existing user.
// All fields are optional, allowing partial updates of user information.
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=1,max=100"` // Updated user name (1-100 characters, optional)
	Email string `json:"email" validate:"omitempty,email"`        // Updated email address (valid email format, optional)
}

// LoginRequest represents the request payload for user authentication.
// It contains the credentials needed to authenticate a user and generate access tokens.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // User's email address for authentication
	Password string `json:"password" validate:"required"`    // User's password for authentication
}

// RegisterRequest represents the request payload for user registration.
// It contains all the information needed to create a new user account.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"` // User's full name (1-100 characters)
	Email    string `json:"email" validate:"required,email"`        // User's email address (must be valid email format)
	Password string `json:"password" validate:"required,min=6"`     // User's password (minimum 6 characters)
}

// AuthResponse represents the response payload for successful authentication.
// It contains the user information and tokens needed for subsequent API requests.
type AuthResponse struct {
	User         *User  `json:"user"`         // Authenticated user's profile information
	AccessToken  string `json:"accessToken"`  // JWT access token for API authentication
	RefreshToken string `json:"refreshToken"` // JWT refresh token for obtaining new access tokens
	ExpiresIn    int64  `json:"expiresIn"`    // Access token expiration time in seconds
}
