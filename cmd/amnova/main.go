package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/am-kenny/amnova/internal/health"
	"github.com/am-kenny/amnova/internal/httpapi"
)

const listenAddr = ":8080"

func main() {
	slog.Info("starting AmNova")

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv, err := httpapi.Listen(listenAddr, health.NewHandler())
	if err != nil {
		slog.Error("server init failed", "err", err)
		os.Exit(1)
	}
	slog.Info("server listening", "addr", srv.Addr())

	if err := srv.Run(ctx); err != nil {
		slog.Error("server stopped with an error", "err", err)
		os.Exit(1)
	}

	slog.Info("Shutting down")
}
