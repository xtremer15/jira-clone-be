# Jira Clone Backend

A scalable, maintainable Go backend for a Jira-like board application built with clean architecture principles.

## Features

- **Authentication & Authorization**: JWT-based authentication with refresh tokens
- **Project Management**: Create, read, update, and delete projects
- **Issue Management**: Full CRUD operations for issues with filtering and sorting
- **User Management**: User registration, authentication, and profile management
- **Rate Limiting**: Built-in rate limiting to prevent abuse
- **Clean Architecture**: Follows SOLID principles with clear separation of concerns
- **JSON Storage**: Mock JSON storage for development (easily migratable to real database)
- **Docker Support**: Complete Docker setup for containerization
- **Health Checks**: Built-in health check endpoints
- **Comprehensive Logging**: Structured JSON logging for monitoring

## Architecture

The application follows clean architecture principles with the following layers:

```
cmd/server/          # Application entry point and server setup
internal/
├── config/          # Configuration management and environment variables
├── routes/          # Centralized route configuration and middleware setup
├── repository/      # Data access layer (interfaces + JSON file implementations)
├── service/         # Business logic layer and validation
├── handler/         # HTTP request/response handlers
└── middleware/      # HTTP middleware (auth, rate limiting, CORS)
types/               # Domain models, DTOs, and enums
pkg/                 # Reusable packages (structured logging)
tests/               # Unit tests and mock implementations
```

## 📁 Project Structure

```
jira-clone-be/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point and server setup
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management and environment variables
│   ├── routes/
│   │   └── routes.go               # Centralized route configuration and middleware setup
│   ├── repository/
│   │   ├── interfaces.go           # Repository interfaces and contracts
│   │   ├── user_repository.go      # User data access layer (JSON file storage)
│   │   ├── project_repository.go  # Project data access layer (JSON file storage)
│   │   └── issue_repository.go    # Issue data access layer (JSON file storage)
│   ├── service/
│   │   ├── user_service.go         # User business logic and validation
│   │   ├── project_service.go      # Project business logic and board management
│   │   ├── issue_service.go        # Issue business logic and workflow
│   │   └── auth_service.go         # Authentication, JWT tokens, and security
│   ├── handler/
│   │   ├── user_handler.go         # User HTTP request/response handlers
│   │   ├── project_handler.go      # Project HTTP request/response handlers
│   │   ├── issue_handler.go        # Issue HTTP request/response handlers
│   │   └── auth_handler.go         # Authentication HTTP handlers (login/register)
│   └── middleware/
│       ├── auth.go                 # JWT authentication middleware
│       └── rate_limit.go           # API rate limiting middleware
├── types/
│   ├── enums.go                    # Domain enums (IssueType, IssueStatus, etc.)
│   ├── user.go                     # User domain model and DTOs
│   ├── project.go                  # Project domain model and DTOs
│   ├── issue.go                    # Issue domain model and DTOs
│   └── board.go                    # Board and Lane domain models
├── pkg/
│   └── logger/
│       └── logger.go               # Structured JSON logging with levels
├── tests/
│   ├── mocks.go                    # Mock implementations for testing
│   ├── user_service_test.go        # User service unit tests
│   └── auth_service_test.go        # Authentication service unit tests
├── scripts/
│   ├── start.sh                    # Linux/macOS startup script
│   └── start.bat                   # Windows startup script
├── Dockerfile                      # Docker container configuration
├── docker-compose.yaml             # Docker Compose for development
├── .env.example                    # Environment variables template
├── .gitignore                      # Git ignore rules
├── go.mod                          # Go module dependencies
├── go.sum                          # Go module checksums
└── README.md                       # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.22+
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd jira-clone-be
```

2. Install dependencies:
```bash
go mod download
```

3. Copy environment variables:
```bash
cp env.example .env
```

4. Update the `.env` file with your configuration:
```bash
# Server Configuration
PORT=8080
LOG_LEVEL=info

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Rate Limiting
RATE_LIMIT=100
```

### Running the Application

#### Option 1: Direct Go execution
```bash
go run cmd/server/main.go
```

#### Option 2: Using Docker
```bash
# Build and run with Docker Compose
docker-compose up --build

# Or build and run with Docker
docker build -t jira-clone-be .
docker run -p 8080:8080 jira-clone-be
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh access token

### Users (Protected)
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user

### Projects (Protected)
- `GET /api/v1/projects` - Get all projects
- `POST /api/v1/projects` - Create project
- `GET /api/v1/projects/{id}` - Get project by ID
- `PUT /api/v1/projects/{id}` - Update project
- `DELETE /api/v1/projects/{id}` - Delete project

### Issues (Protected)
- `GET /api/v1/issues` - Get all issues
- `POST /api/v1/issues` - Create issue
- `GET /api/v1/issues/{id}` - Get issue by ID
- `PUT /api/v1/issues/{id}` - Update issue
- `DELETE /api/v1/issues/{id}` - Delete issue

### Health Check
- `GET /healthz` - Health check endpoint

## Data Models

### User
```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "avatarUrl": "string",
  "createdAt": "ISO8601 string",
  "updatedAt": "ISO8601 string",
  "issueIds": ["array of string"]
}
```

### Project
```json
{
  "id": "string",
  "name": "string",
  "url": "string",
  "description": "string",
  "category": "enum: [Software, Marketing, Business]",
  "createdAt": "ISO8601 string",
  "updatedAt": "ISO8601 string",
  "issues": ["array of Issue objects"],
  "users": ["array of User objects"]
}
```

### Issue
```json
{
  "id": "string",
  "priority": "enum: [LOWEST, LOW, MEDIUM, HIGH, HIGHEST]",
  "status": "enum: [BACKLOG, SELECTED, IN_PROGRESS, DONE]",
  "title": "string",
  "description": "string",
  "type": "enum: [BUG, STORY, TASK]",
  "icon": "string",
  "assignee": "string",
  "reporter": "string",
  "iconUrl": "string",
  "displayId": "string",
  "projectId": "string",
  "createdAt": "ISO8601 string",
  "updatedAt": "ISO8601 string"
}
```

## Development

### Running Tests
```bash
go test ./...
```

### Code Quality
```bash
# Format code
gofmt -s -w .

# Lint code
golint ./...

# Vet code
go vet ./...
```

### Building for Production
```bash
# Build binary
go build -o bin/jira-clone-be cmd/server/main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/jira-clone-be-linux cmd/server/main.go
```

## Configuration

The application supports the following environment variables:

- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (debug, info, warn, error, fatal)
- `JWT_SECRET`: JWT signing secret (change in production!)
- `RATE_LIMIT`: Rate limit per minute (default: 100)

## Security Features

- JWT-based authentication with refresh tokens
- Input validation and sanitization
- Rate limiting to prevent abuse
- CORS configuration for frontend integration
- Secure password hashing (bcrypt)

## Monitoring

- Health check endpoint at `/healthz`
- Structured JSON logging
- Request ID tracking
- Error recovery middleware

## Future Enhancements

- Database integration (PostgreSQL/MongoDB)
- Redis caching
- GraphQL API
- WebSocket support for real-time updates
- File upload support
- Advanced search and filtering
- Audit logging
- Metrics and monitoring integration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Adding a new route (step-by-step tutorial)

This section shows the minimal, recommended steps to add a new API route to this project so it becomes available and eligible for calling. The project follows a layered structure (handler → service → repository) and a new route usually touches files in each layer plus the shared `types` and `api/openapi.yaml`.

Short checklist:
- Add or update a DTO/domain type in `types/` (if needed)
- Add repository interface method(s) (if data persistence is required)
- Implement repository logic (file-backed implementation under `internal/repository/`)
- Add service method(s) under `internal/service/`
- Add handler method(s) under `internal/handler/`
- Wire repository → service → handler in `internal/routes/routes.go`
- Register the HTTP route in `Routes.SetupRoutes()` (public or protected)
- Update `api/openapi.yaml` for docs and clients
- Add unit tests under `tests/` and update `tests/mocks.go` if necessary

Detailed steps with examples

1) Define request/response DTOs (optional)

- If your route accepts or returns structured JSON that doesn't match existing domain models, add a new file in `types/` (for example `types/foo.go`) and declare the request/response structs and validation tags using `github.com/go-playground/validator/v10` tags.

Example: add `types.Foo` and `types.CreateFooRequest`.

2) Repository (persistence) — interface and implementation (optional)

- If the route requires storing or reading from persistence, add method signatures to the repository interface in `internal/repository/interfaces.go` (for the appropriate repository: UserRepository, ProjectRepository, or IssueRepository). For a new resource type you can create a new repository interface and implementation file under `internal/repository/`.

- Implement the file-backed repository under `internal/repository/` (follow existing patterns in `user_repository.go` / `issue_repository.go`). Make sure to:
  - keep a map in memory with a `sync.RWMutex`
  - load data from `data/<resource>.json` in `loadData()`
  - call `saveData()` after writes
  - return clear errors for not-found cases

3) Service layer

- Add a new service file under `internal/service/` (e.g., `foo_service.go`). Implement a constructor `NewFooService(repo repository.FooRepository, logger *logger.Logger) *FooService` and the business methods your handler will call. Keep services thin: validate inputs, perform domain logic, call repository methods, and log appropriately.

4) Handler layer

- Add a handler file under `internal/handler/` (e.g., `foo_handler.go`) with a constructor `NewFooHandler(fooService *service.FooService, logger *logger.Logger) *FooHandler` and HTTP methods matching the route(s) (e.g., `CreateFoo`, `GetFoo`). Use `encoding/json` to decode request bodies, `validator.New()` to validate structs, and `render.JSON` to reply.

Handler responsibilities:
- Decode request
- Validate input (`validator`)
- Call service
- Map service errors to HTTP status codes
- Return JSON with `render.JSON` or appropriate status code

5) Wire everything in the composition root

- Open `internal/routes/routes.go`. In `NewRoutes` add construction of the new repository, service and handler alongside the existing ones. Example pattern:

  // Initialize repositories
  fooRepo := repository.NewFooRepository()

  // Initialize services
  fooService := service.NewFooService(fooRepo, logger)

  // Initialize handlers
  fooHandler := handler.NewFooHandler(fooService, logger)

Add these to the `Routes` struct if you want them available on the receiver (optional but recommended for consistency).

6) Register the HTTP route(s)

- Edit `Routes.SetupRoutes()` to register the new endpoint under the desired group (public or protected). For protected routes ensure the route is inside the `protectedRouter` group which uses the JWT `AuthMiddleware`.

Example: add a protected POST route for creating a Foo under `/api/v1`:

  protectedRouter.Post("/foo", r.fooHandler.CreateFoo)

If it's public, register it on the outer `apiRouter` before the protected group.

7) Update OpenAPI documentation

- Update `api/openapi.yaml` to include the new path, request schema and response schema so generated clients and docs stay accurate. Follow the existing YAML structure used for other endpoints.

8) Tests and mocks

- Add unit tests under `tests/` for your service and handler. If you added a new repository interface, add a mock implementation in `tests/mocks.go` or create a new file under `tests/` and use `testify/mock` as existing tests do.

- Run tests:

```bash
go test ./...
```

9) Formatting and static checks

- Run format and vet tools before committing:

```bash
gofmt -s -w .
go vet ./...
```

10) Example end-to-end (create a simple GET route)

- Suppose you want a simple, protected GET `/api/v1/ping` that returns `{ "pong": true }`.

  - You can add a lightweight handler method directly: create `internal/handler/ping_handler.go` with a `Ping` method that writes `render.JSON(w, r, map[string]bool{"pong": true})`. Add a `NewPingHandler` constructor.
  - In `routes.NewRoutes` create a ping handler: `pingHandler := handler.NewPingHandler(logger)` and keep it on the returned `Routes` struct.
  - In `Routes.SetupRoutes()` inside the protected group add:

    protectedRouter.Get("/ping", r.pingHandler.Ping)

  - Optionally add the path to `api/openapi.yaml`.

That’s it — after wiring and building, the new route will be available on the running server.

Notes and best practices
- Keep handler logic minimal; put business rules in services.
- Keep repository interfaces stable: add methods carefully and keep semantics clear (e.g., Create/GetByID/GetAll/Update/Delete).
- Add validation tags to `types` structs and validate in handlers.
- Protect routes with `AuthMiddleware` unless the route must be public.
- For production or multi-instance deployments, replace the in-memory rate limiter and file-backed repo with centralized services (Redis, Postgres).

If you want, I can add a small example route implementation (handler + route registration + a test) to this repository to demonstrate all steps concretely — tell me which endpoint you'd like (e.g., `GET /api/v1/ping` or `POST /api/v1/foos`).
