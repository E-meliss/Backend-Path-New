package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimitConfig struct {
	RPS   int
	Burst int
}

type bucket struct {
	tokens float64
	last   time.Time
}

func RateLimit(cfg RateLimitConfig) Middleware {
	var mu sync.Mutex
	buckets := map[string]*bucket{}

	refillRate := float64(cfg.RPS)
	maxTokens := float64(cfg.Burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)

			mu.Lock()
			b, ok := buckets[ip]
			if !ok {
				b = &bucket{tokens: maxTokens, last: time.Now()}
				buckets[ip] = b
			}
			now := time.Now()
			elapsed := now.Sub(b.last).Seconds()
			b.tokens = min(maxTokens, b.tokens+elapsed*refillRate)
			b.last = now

			if b.tokens < 1 {
				mu.Unlock()
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			b.tokens -= 1
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
