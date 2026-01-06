package http

import (
	"log/slog"
	nethttp "net/http"
	"time"

	"github.com/E-meliss/wallet-service/internal/cache"
	"github.com/E-meliss/wallet-service/internal/config"
	"github.com/E-meliss/wallet-service/internal/eventstore"
	"github.com/E-meliss/wallet-service/internal/http/handlers"
	"github.com/E-meliss/wallet-service/internal/http/middleware"
)

type Deps struct {
	Cfg   config.Config
	Log   *slog.Logger
	Cache cache.Cache
	ES    eventstore.Store
}

func NewServer(addr string, deps Deps) *nethttp.Server {
	r := NewRouter()

	r.Use(
		middleware.RequestID(),
		middleware.Recover(deps.Log),
		middleware.SecurityHeaders(),
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins: []string{"*"}, // tighten for prod
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Authorization", "Content-Type", "X-Request-Id"},
		}),
		middleware.RateLimit(middleware.RateLimitConfig{RPS: 10, Burst: 20}),
		middleware.CircuitBreaker(middleware.CircuitBreakerConfig{FailureThreshold: 20, OpenFor: 10 * time.Second}),
		middleware.Logging(deps.Log),
		middleware.Metrics(),
	)

	r.Handle(nethttp.MethodGet, "/health", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	r.Handle(nethttp.MethodGet, "/metrics", middleware.MetricsHandler())

	authMW := middleware.AuthStub()
	requireRoleGen := func(role string) func(nethttp.Handler) nethttp.Handler { return middleware.RequireRole(role) }

	authH := handlers.NewAuthHandler()
	usersH := handlers.NewUsersHandler()
	txH := handlers.NewTransactionsHandler()
	balH := handlers.NewBalancesHandler()

	api := NewV1Routes(deps)
	api.Register(r,
		authH,
		usersH,
		txH,
		balH,
		authMW,
		requireRoleGen,
	)

	return &nethttp.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
