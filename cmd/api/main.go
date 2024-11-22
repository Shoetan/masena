package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/masena-dev/bookstore-api/internal/api"
	"github.com/masena-dev/bookstore-api/internal/db"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config holds the configuration settings for the application.
type config struct {
	Port   int    `envconfig:"PORT" default:"4000"`
	DB_URL string `envconfig:"DATABASE_URL" required:"true"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	loggerOpts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, loggerOpts))

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DB_URL)
	if err != nil {
		logger.Error("error setting up db", "error", err)
		os.Exit(1)
	}

	defer pool.Close()

	srv := api.NewServer(logger, db.New(pool))

	if err = serve(cfg, logger, srv); err != nil {
		logger.Error("error starting server", "error", err)
		os.Exit(1)
	}
}

// loadConfig loads the configuration from environment variables using envconfig.
func loadConfig() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// serve starts the HTTP server and handles graceful shutdown.
func serve(cfg config, logger *slog.Logger, srv *api.Server) error {
	handler := srv.NewRouter()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Info("shutdown", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			shutdownError <- err
		}

		logger.Info("completing background tasks")
		shutdownError <- nil
	}()

	logger.Info("starting server", "addr", httpServer.Addr)

	err := httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownError
}
