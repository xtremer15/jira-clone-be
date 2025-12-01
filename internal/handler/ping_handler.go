// filepath: g:\Golang\jira-clone-be\internal\handler\ping_handler.go
package handler

import (
    "net/http"

    "jira-clone-be/pkg/logger"

    "github.com/go-chi/render"
)

// PingHandler handles simple ping endpoints used for smoke checks
type PingHandler struct {
    logger *logger.Logger
}

// NewPingHandler creates a new PingHandler
func NewPingHandler(logger *logger.Logger) *PingHandler {
    return &PingHandler{logger: logger}
}

// Ping returns a simple pong response
func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
    // Optionally you can access userID from context when this route is protected
    // userID := r.Context().Value("userID")
    render.JSON(w, r, map[string]bool{"pong": true})
}
