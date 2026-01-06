package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/E-meliss/wallet-service/internal/config"
	apphttp "github.com/E-meliss/wallet-service/internal/http"
)

func New(addr string, cfg config.Config, log *slog.Logger) *http.Server {
	return apphttp.NewServer(addr, apphttp.Deps{Cfg: cfg, Log: log})
}
