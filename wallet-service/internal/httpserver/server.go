package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	apphttp "github.com/E-meliss/wallet-service/internal/http"
)

// New wires the custom router + middleware stack from internal/http.
func New(addr string, log *slog.Logger, pool *pgxpool.Pool) *http.Server {
	return apphttp.NewServer(addr, apphttp.Deps{
		Log: log,
		DB:  pool,
	})
}
