package middleware

import (
	"net/http"
	"sync"
	"time"
)

type CircuitBreakerConfig struct {
	FailureThreshold int
	OpenFor          time.Duration
}

type breakerState struct {
	mu        sync.Mutex
	failures  int
	openUntil time.Time
}

func CircuitBreaker(cfg CircuitBreakerConfig) Middleware {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 20
	}
	if cfg.OpenFor <= 0 {
		cfg.OpenFor = 10 * time.Second
	}

	state := &breakerState{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state.mu.Lock()
			open := time.Now().Before(state.openUntil)
			state.mu.Unlock()
			if open {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("circuit breaker open"))
				return
			}

			rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			if rw.status >= 500 {
				state.mu.Lock()
				state.failures++
				if state.failures >= cfg.FailureThreshold {
					state.openUntil = time.Now().Add(cfg.OpenFor)
					state.failures = 0
				}
				state.mu.Unlock()
				return
			}

			state.mu.Lock()
			if state.failures > 0 {
				state.failures--
			}
			state.mu.Unlock()
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
