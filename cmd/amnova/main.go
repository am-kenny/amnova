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
	"github.com/am-kenny/amnova/internal/store"
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

	st, err := store.Open(ctx, cfg.Storage.DataDir)
	if err != nil {
		slog.Error("store open failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := st.Close(); err != nil {
			slog.Error("store close failed", "err", err)
		}
	}()
	slog.Info("store opened", "data_dir", cfg.Storage.DataDir)

	if err := st.Migrate(ctx); err != nil {
		slog.Error("migration failed", "err", err)
		os.Exit(1)
	}

	srv, err := httpapi.Listen(cfg.Server.ListenAddr, health.NewHandler(st))
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
