// filepath: g:\Golang\jira-clone-be\tests\ping_handler_test.go
package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "jira-clone-be/internal/handler"
    "jira-clone-be/pkg/logger"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
)

func TestPingHandler_Ping(t *testing.T) {
    logger := logger.New("info")
    pingHandler := handler.NewPingHandler(logger)

    r := chi.NewRouter()
    r.Get("/ping", pingHandler.Ping)

    req := httptest.NewRequest(http.MethodGet, "/ping", nil)
    rr := httptest.NewRecorder()

    r.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    assert.JSONEq(t, `{"pong": true}`, rr.Body.String())
}
