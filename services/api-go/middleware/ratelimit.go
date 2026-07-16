package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/oni1997/onentry/services/api-go/models"
)

type RateLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
	rpm     int
	window  time.Duration
}

type client struct {
	requests int
	lastReset time.Time
}

func NewRateLimiter(rpm int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rpm:     rpm,
		window:  time.Minute,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		if !rl.allow(ip) {
			w.Header().Set("Retry-After", "60")
			models.JSONError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &client{lastReset: now}
		c = rl.clients[ip]
	}

	if now.Sub(c.lastReset) >= rl.window {
		c.requests = 0
		c.lastReset = now
	}

	if c.requests >= rl.rpm {
		return false
	}

	c.requests++
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, c := range rl.clients {
			if now.Sub(c.lastReset) > rl.window*2 {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}
