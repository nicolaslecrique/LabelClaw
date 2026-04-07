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
	"github.com/nicolaslecrique/LabelClaw/backend/internal/config"
	"github.com/nicolaslecrique/LabelClaw/backend/internal/configuration"
	"github.com/nicolaslecrique/LabelClaw/backend/internal/storage"
)

func main() {
	cfg := config.Load()

	store := storage.NewFileStore(cfg.ConfigurationPath)
	service := configuration.NewService(store)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           api.NewHandler(service),
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

	log.Printf("labelclaw backend listening on %s", cfg.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
