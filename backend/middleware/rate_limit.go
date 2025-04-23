// middleware/rate_limit.go
package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.limiters[ip]; exists {
		return limiter
	}

	// Allow 60 requests per minute
	limiter := rate.NewLimiter(rate.Limit(1), 60)
	rl.limiters[ip] = limiter
	return limiter
}

func RateLimit(next http.Handler) http.Handler {
	rl := NewRateLimiter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := rl.getLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
