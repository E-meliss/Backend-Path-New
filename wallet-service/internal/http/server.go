package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/E-meliss/wallet-service/internal/http/middleware"
)

type Deps struct {
	Log *slog.Logger
	DB  *pgxpool.Pool
}

func NewServer(addr string, deps Deps) *http.Server {
	r := NewRouter()

	r.Use(
		middleware.RequestID(),
		middleware.Recover(deps.Log),
		middleware.SecurityHeaders(),
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins: []string{"*"}, // prod’da daralt
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Authorization", "Content-Type", "X-Request-Id"},
		}),
		middleware.RateLimit(middleware.RateLimitConfig{
			RPS:   10,
			Burst: 20,
		}),
		middleware.Logging(deps.Log),
		middleware.Metrics(),
	)

	r.Handle(http.MethodGet, "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	api := NewV1Routes(deps)
	api.Register(r)

	return &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
