package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nicolas/labelclaw/backend/internal/api"
	"github.com/nicolas/labelclaw/backend/internal/configuration"
	"github.com/nicolas/labelclaw/backend/internal/llm"
	"github.com/nicolas/labelclaw/backend/internal/static"
)

func main() {
	port := envOrDefault("PORT", "8080")
	dataDir := envOrDefault("DATA_DIR", "./data")
	storePath := filepath.Join(dataDir, "active-config.json")

	store := configuration.NewFileStore(storePath)
	staticHandler, err := static.NewHandler()
	if err != nil {
		log.Fatalf("create static handler: %v", err)
	}

	var generator llm.Client
	if config, ok := llm.LoadOpenAICompatibleConfigFromEnv(); ok {
		generator = llm.NewOpenAICompatibleClient(config)
	} else {
		generator = llm.NewUnavailableClient("LLM provider is not configured. Set LLM_BASE_URL, LLM_API_KEY, and LLM_MODEL.")
	}

	server := api.NewServer(api.Dependencies{
		Store:         store,
		Generator:     generator,
		StaticHandler: staticHandler,
		AllowedOrigin: envOrDefault("CORS_ALLOWED_ORIGIN", "http://127.0.0.1:5173"),
	})

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           server,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("labelclaw listening on http://127.0.0.1:%s", port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}

func envOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

