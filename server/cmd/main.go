package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VaraPrasad27/shop/server/internal/config"
	"github.com/VaraPrasad27/shop/server/internal/db"
	"github.com/VaraPrasad27/shop/server/internal/routes"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("startup: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	dbpool, err := db.Connect(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer dbpool.Close()

	r := routes.New(dbpool)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Run the server in a goroutine so we can wait for a signal.
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	// Block until we get a signal or the server fails.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server: %w", err)
	case sig := <-stop:
		log.Printf("received %s, shutting down", sig)
	}

	// Give in-flight requests up to 15s to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	// Drain any server-error signal that arrived during shutdown.
	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("server: %w", err)
		}
	case <-time.After(5 * time.Second):
	}

	return nil
}
