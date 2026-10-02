package main

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/datasektionen/OG-template/internal/config"
	"github.com/datasektionen/OG-template/internal/database"
	"github.com/datasektionen/OG-template/internal/web"
)

//go:embed schema.sql
var schemaSQL string

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	env := config.GetEnv()
	slog.Info("Starting OG-template service", "port", env.ServerPort, "public_url", env.ServerURL)

	// Initialize database and apply schema if DB_CONN is configured
	db, err := database.Init(env.DBConn, schemaSQL)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	if db != nil {
		defer db.Close()
	}

	handler, err := web.NewHandler(db, env)
	if err != nil {
		slog.Error("Failed to mount web routes", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:         "0.0.0.0:" + env.ServerPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown on SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Server listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received, gracefully draining connections...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed, forcing close", "error", err)
	} else {
		slog.Info("Server exited cleanly")
	}
}
