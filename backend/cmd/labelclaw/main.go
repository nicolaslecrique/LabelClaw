package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicolaslecrique/LabelClaw/backend/internal/api"
	"github.com/nicolaslecrique/LabelClaw/backend/internal/storage"
)

func main() {
	addr := getenv("LABELCLAW_ADDR", "127.0.0.1:8080")
	configPath := getenv("LABELCLAW_CONFIG_PATH", "data/active-config.json")

	server := &http.Server{
		Addr:              addr,
		Handler:           api.NewHandler(storage.NewFileStore(configPath)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownContext.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	}()

	log.Printf("labelclaw backend listening on %s", addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
