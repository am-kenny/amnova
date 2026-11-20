package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/am-kenny/amnova/internal/config"
	"github.com/am-kenny/amnova/internal/health"
	"github.com/am-kenny/amnova/internal/httpapi"
)

func main() {
	// Log as JSON from the start
	slog.SetDefault(newLogger(slog.LevelInfo))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}
	slog.SetDefault(newLogger(cfg.Log.Level))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv, err := httpapi.Listen(cfg.Server.ListenAddr, health.NewHandler())
	if err != nil {
		slog.Error("server init failed", "err", err)
		os.Exit(1)
	}
	slog.Info("server listening", "addr", srv.Addr())

	if err := srv.Run(ctx); err != nil {
		slog.Error("server stopped with an error", "err", err)
		os.Exit(1)
	}

	slog.Info("shutting down")
}

func newLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
