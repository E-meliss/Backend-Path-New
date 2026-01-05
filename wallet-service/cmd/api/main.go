package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/E-meliss/wallet-service/internal/app"
	"github.com/E-meliss/wallet-service/internal/config"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = a.Shutdown(shutdownCtx)
	}()

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
