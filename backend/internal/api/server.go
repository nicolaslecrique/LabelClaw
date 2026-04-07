package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nicolaslecrique/LabelClaw/backend/internal/configuration"
	"github.com/nicolaslecrique/LabelClaw/backend/internal/storage"
)

type Store interface {
	Load() (configuration.ActiveConfiguration, error)
	Save(configuration.ActiveConfiguration) error
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(store Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/configuration/current", handleGetCurrentConfiguration(store))
	mux.HandleFunc("PUT /api/configuration/current", handlePutCurrentConfiguration(store))
	mux.HandleFunc("POST /api/configuration/generate", handleGenerateConfiguration)

	return mux
}

func handleHealth(responseWriter http.ResponseWriter, _ *http.Request) {
	writeJSON(responseWriter, http.StatusOK, map[string]string{"status": "ok"})
}

func handleGetCurrentConfiguration(store Store) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, _ *http.Request) {
		current, err := store.Load()
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeJSON(responseWriter, http.StatusNotFound, errorResponse{Error: err.Error()})
				return
			}

			writeJSON(responseWriter, http.StatusInternalServerError, errorResponse{Error: "failed to load active configuration"})
			return
		}

		writeJSON(responseWriter, http.StatusOK, current)
	}
}

func handlePutCurrentConfiguration(store Store) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		var current configuration.ActiveConfiguration
		if err := decodeJSON(request, &current); err != nil {
			writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		if err := current.Validate(); err != nil {
			writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		if err := store.Save(current); err != nil {
			writeJSON(responseWriter, http.StatusInternalServerError, errorResponse{Error: "failed to save active configuration"})
			return
		}

		responseWriter.WriteHeader(http.StatusNoContent)
	}
}

func handleGenerateConfiguration(responseWriter http.ResponseWriter, request *http.Request) {
	var generateRequest configuration.GenerateRequest
	if err := decodeJSON(request, &generateRequest); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := generateRequest.Validate(); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(responseWriter, http.StatusNotImplemented, errorResponse{Error: configuration.ErrGenerateNotImplemented.Error()})
}

func decodeJSON(request *http.Request, destination any) error {
	defer request.Body.Close()

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	return nil
}

func writeJSON(responseWriter http.ResponseWriter, status int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)

	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
