package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nicolas/labelclaw/backend/internal/configuration"
	"github.com/nicolas/labelclaw/backend/internal/llm"
)

type Dependencies struct {
	Store         configuration.Store
	Generator     llm.Client
	StaticHandler http.Handler
	AllowedOrigin string
}

type Server struct {
	store         configuration.Store
	generator     llm.Client
	staticHandler http.Handler
	allowedOrigin string
}

func NewServer(dependencies Dependencies) http.Handler {
	server := &Server{
		store:         dependencies.Store,
		generator:     dependencies.Generator,
		staticHandler: dependencies.StaticHandler,
		allowedOrigin: dependencies.AllowedOrigin,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", server.handleHealth)
	mux.HandleFunc("GET /api/configuration/current", server.handleLoadCurrentConfiguration)
	mux.HandleFunc("PUT /api/configuration/current", server.handleSaveCurrentConfiguration)
	mux.HandleFunc("POST /api/configuration/generate", server.handleGenerateConfiguration)
	mux.Handle("/", server.staticHandler)

	return server.withCORS(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleLoadCurrentConfiguration(w http.ResponseWriter, _ *http.Request) {
	config, err := s.store.Load()
	if errors.Is(err, configuration.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "No saved configuration found."})
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to load the active configuration.")
		return
	}

	writeJSON(w, http.StatusOK, config)
}

func (s *Server) handleSaveCurrentConfiguration(w http.ResponseWriter, r *http.Request) {
	var payload configuration.SavedConfiguration
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Request body must be valid JSON.")
		return
	}

	if err := validateSavedConfiguration(payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	payload.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.store.Save(payload); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to save the active configuration.")
		return
	}

	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleGenerateConfiguration(w http.ResponseWriter, r *http.Request) {
	var payload configuration.GenerateInput
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Request body must be valid JSON.")
		return
	}

	if err := validateGenerateInput(payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := s.generator.GenerateLabelingPanel(r.Context(), llm.GenerateRequest{
		SampleSchema: payload.SampleSchema,
		LabelSchema:  payload.LabelSchema,
		UIPrompt:     payload.UIPrompt,
	})
	if err != nil {
		if errors.Is(err, llm.ErrProviderUnavailable) {
			writeError(w, http.StatusServiceUnavailable, err.Error())
			return
		}

		writeError(w, http.StatusBadGateway, fmt.Sprintf("LLM generation failed: %v", err))
		return
	}

	if err := configuration.ValidateGeneratedComponentSource(response.ComponentSource); err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	if err := configuration.ValidateDataAgainstSchema(payload.SampleSchema, response.SampleData); err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			if s.allowedOrigin == "*" || origin == s.allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateGenerateInput(input configuration.GenerateInput) error {
	if err := configuration.ValidateSchemaJSON(input.SampleSchema); err != nil {
		return fmt.Errorf("sampleSchema: %w", err)
	}

	if err := configuration.ValidateSchemaJSON(input.LabelSchema); err != nil {
		return fmt.Errorf("labelSchema: %w", err)
	}

	if strings.TrimSpace(input.UIPrompt) == "" {
		return errors.New("uiPrompt must not be empty")
	}

	return nil
}

func validateSavedConfiguration(config configuration.SavedConfiguration) error {
	if err := configuration.ValidateSchemaJSON(config.SampleSchema); err != nil {
		return fmt.Errorf("sampleSchema: %w", err)
	}

	if err := configuration.ValidateSchemaJSON(config.LabelSchema); err != nil {
		return fmt.Errorf("labelSchema: %w", err)
	}

	if strings.TrimSpace(config.UIPrompt) == "" {
		return errors.New("uiPrompt must not be empty")
	}

	if err := configuration.ValidateDataAgainstSchema(config.SampleSchema, config.SampleData); err != nil {
		return fmt.Errorf("sampleData: %w", err)
	}

	if err := configuration.ValidateGeneratedComponentSource(config.ComponentSource); err != nil {
		return fmt.Errorf("componentSource: %w", err)
	}

	return nil
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"message": message})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

