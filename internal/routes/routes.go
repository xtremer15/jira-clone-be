// Package routes provides centralized route configuration and middleware setup for the Jira Clone Backend.
// It handles the initialization of all dependencies and configures HTTP routes with proper middleware.
package routes

import (
	"net/http"
	"time"

	"jira-clone-be/internal/config"
	"jira-clone-be/internal/handler"
	authMiddleware "jira-clone-be/internal/middleware"
	"jira-clone-be/internal/repository"
	"jira-clone-be/internal/service"
	"jira-clone-be/pkg/logger"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Routes holds all the route dependencies and provides centralized route configuration.
// It manages the initialization of repositories, services, handlers, and middleware.
type Routes struct {
	config *config.Config
	logger *logger.Logger

	// Repositories
	userRepo    repository.UserRepository
	projectRepo repository.ProjectRepository
	issueRepo   repository.IssueRepository

	// Services
	userService    *service.UserService
	projectService *service.ProjectService
	issueService   *service.IssueService
	authService    *service.AuthService

	// Handlers
	userHandler    *handler.UserHandler
	projectHandler *handler.ProjectHandler
	issueHandler   *handler.IssueHandler
	authHandler    *handler.AuthHandler
	pingHandler    *handler.PingHandler
}

// NewRoutes creates a new routes instance with all dependencies initialized.
// It sets up repositories, services, handlers, and returns a configured Routes instance.
//
// Parameters:
//   - cfg: Application configuration containing server settings and JWT secrets
//   - logger: Structured logger instance for logging throughout the application
//
// Returns a new Routes instance with all dependencies properly initialized.
func NewRoutes(cfg *config.Config, logger *logger.Logger) *Routes {
	// Initialize repositories
	userRepo := repository.NewUserRepository()
	projectRepo := repository.NewProjectRepository()
	issueRepo := repository.NewIssueRepository()

	// Initialize services
	userService := service.NewUserService(userRepo, logger)
	projectService := service.NewProjectService(projectRepo, logger)
	issueService := service.NewIssueService(issueRepo, logger)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, logger)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, logger)
	projectHandler := handler.NewProjectHandler(projectService, logger)
	issueHandler := handler.NewIssueHandler(issueService, logger)
	authHandler := handler.NewAuthHandler(authService, logger)
	pingHandler := handler.NewPingHandler(logger)

	return &Routes{
		config: cfg,
		logger: logger,

		// Repositories
		userRepo:    userRepo,
		projectRepo: projectRepo,
		issueRepo:   issueRepo,

		// Services
		userService:    userService,
		projectService: projectService,
		issueService:   issueService,
		authService:    authService,

		// Handlers
		userHandler:    userHandler,
		projectHandler: projectHandler,
		issueHandler:   issueHandler,
		authHandler:    authHandler,
		pingHandler:    pingHandler,
	}
}

// SetupRoutes configures and returns all routes with middleware.
// It sets up public routes (health check, authentication) and protected routes
// (user, project, issue management) with proper middleware configuration.
//
// Returns a configured Chi router with all routes and middleware applied.
// SetupRoutes configures and returns all routes with middleware.
func (r *Routes) SetupRoutes() *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	r.setupMiddleware(router)

	// Health check (nu ține de /api/v1)
	router.Get("/healthz", r.healthCheck)

	// API v1 routes
	router.Route("/api/v1", func(apiRouter chi.Router) {
		// --- Public routes ---
		apiRouter.Post("/auth/login", r.authHandler.Login)
		apiRouter.Post("/auth/register", r.authHandler.Register)
		apiRouter.Post("/auth/refresh", r.authHandler.RefreshToken)

		// --- Protected routes ---
		apiRouter.Group(func(protectedRouter chi.Router) {
			// Auth middleware
			protectedRouter.Use(authMiddleware.AuthMiddleware(r.authService))

			// User routes
			protectedRouter.Get("/users", r.userHandler.GetUsers)
			protectedRouter.Get("/users/{id}", r.userHandler.GetUser)
			protectedRouter.Put("/users/{id}", r.userHandler.UpdateUser)

			// Project routes
			protectedRouter.Get("/projects", r.projectHandler.GetProjects)
			protectedRouter.Post("/projects", r.projectHandler.CreateProject)
			protectedRouter.Get("/projects/{id}", r.projectHandler.GetProject)
			protectedRouter.Put("/projects/{id}", r.projectHandler.UpdateProject)
			protectedRouter.Delete("/projects/{id}", r.projectHandler.DeleteProject)

			// Issue routes
			protectedRouter.Get("/issues", r.issueHandler.GetIssues)
			protectedRouter.Post("/issues", r.issueHandler.CreateIssue)
			protectedRouter.Get("/issues/{id}", r.issueHandler.GetIssue)
			protectedRouter.Put("/issues/{id}", r.issueHandler.UpdateIssue)
			protectedRouter.Delete("/issues/{id}", r.issueHandler.DeleteIssue)

			// Ping (protected example)
			protectedRouter.Get("/ping", r.pingHandler.Ping)
		})
	})

	return router
}

// setupMiddleware configures all middleware for the router.
// It sets up basic middleware (request ID, real IP, logging, recovery, timeout),
// CORS configuration, and rate limiting.
//
// Parameters:
//   - router: The Chi router to configure with middleware
func (r *Routes) setupMiddleware(router *chi.Mux) {
	// Basic middleware
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(chiMiddleware.Timeout(60 * time.Second))

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://angular-jira-clone-alpha.vercel.app", "http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting
	router.Use(authMiddleware.RateLimit(r.config.RateLimit))
}

// setupPublicRoutes configures public routes that do not require authentication.
// These include health check endpoints and authentication endpoints (login, register, refresh).
//
// Parameters:
//   - router: The Chi router to configure with public routes
// func (r *Routes) setupPublicRoutes(router *chi.Mux) {
// 	// Health check
// 	router.Get("/healthz", r.healthCheck)

// 	// API v1 public routes
// 	router.Route("/api/v1", func(apiRouter chi.Router) {
// 		// Authentication routes
// 		apiRouter.Post("/auth/login", r.authHandler.Login)
// 		apiRouter.Post("/auth/register", r.authHandler.Register)
// 		apiRouter.Post("/auth/refresh", r.authHandler.RefreshToken)
// 	})
// }

// setupProtectedRoutes configures protected routes that require authentication.
// These include user management, project management, and issue management endpoints.
// All routes in this group are protected by JWT authentication middleware.
//
// Parameters:
//   - router: The Chi router to configure with protected routes
// func (r *Routes) setupProtectedRoutes(router *chi.Mux) {
// 	router.Route("/api/v1", func(apiRouter chi.Router) {
// 		// Protected routes group
// 		apiRouter.Group(func(protectedRouter chi.Router) {
// 			protectedRouter.Use(authMiddleware.AuthMiddleware(r.authService))

// 			// User routes
// 			protectedRouter.Get("/users", r.userHandler.GetUsers)
// 			protectedRouter.Get("/users/{id}", r.userHandler.GetUser)
// 			protectedRouter.Put("/users/{id}", r.userHandler.UpdateUser)

// 			// Project routes
// 			protectedRouter.Get("/projects", r.projectHandler.GetProjects)
// 			protectedRouter.Post("/projects", r.projectHandler.CreateProject)
// 			protectedRouter.Get("/projects/{id}", r.projectHandler.GetProject)
// 			protectedRouter.Put("/projects/{id}", r.projectHandler.UpdateProject)
// 			protectedRouter.Delete("/projects/{id}", r.projectHandler.DeleteProject)

// 			// Issue routes
// 			protectedRouter.Get("/issues", r.issueHandler.GetIssues)
// 			protectedRouter.Post("/issues", r.issueHandler.CreateIssue)
// 			protectedRouter.Get("/issues/{id}", r.issueHandler.GetIssue)
// 			protectedRouter.Put("/issues/{id}", r.issueHandler.UpdateIssue)
// 			protectedRouter.Delete("/issues/{id}", r.issueHandler.DeleteIssue)
// 		})
// 	})
// }

// healthCheck handles the health check endpoint.
// It returns a simple "OK" response to indicate the server is running and healthy.
// This endpoint is used by load balancers and monitoring systems to check server health.
//
// Parameters:
//   - w: HTTP response writer
//   - req: HTTP request
func (r *Routes) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
