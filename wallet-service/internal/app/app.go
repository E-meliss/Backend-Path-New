package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/E-meliss/wallet-service/internal/config"
	"github.com/E-meliss/wallet-service/internal/db"
	"github.com/E-meliss/wallet-service/internal/httpserver"
	"github.com/E-meliss/wallet-service/internal/logger"
)

type App struct {
	cfg    config.Config
	log    *slog.Logger
	dbPool *pgxpool.Pool
	srv    *http.Server
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	logg := logger.New(cfg.LogLevel)

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	srv := httpserver.New(cfg.HTTPAddr, logg, pool)

	return &App{
		cfg:    cfg,
		log:    logg,
		dbPool: pool,
		srv:    srv,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("server starting", "addr", a.cfg.HTTPAddr)
	return a.srv.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("shutdown started")
	if a.dbPool != nil {
		a.dbPool.Close()
	}
	if a.srv != nil {
		return a.srv.Shutdown(ctx)
	}
	return nil
}
