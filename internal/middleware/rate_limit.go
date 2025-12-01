package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter represents a simple rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rateLimit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    rateLimit,
		window:   time.Minute, // 1 minute window
	}
}

// IsAllowed checks if a request is allowed for the given key
func (rl *RateLimiter) IsAllowed(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old requests
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	// Check if we're under the limit
	requests := rl.requests[key]
	if len(requests) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[key] = append(requests, now)
	return true
}

// RateLimit middleware limits requests per IP
func RateLimit(rateLimit int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rateLimit)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = forwarded
			}

			// Check if request is allowed
			if !limiter.IsAllowed(clientIP) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
