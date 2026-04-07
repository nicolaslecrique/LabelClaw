package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nicolaslecrique/LabelClaw/backend/internal/configuration"
	"github.com/nicolaslecrique/LabelClaw/backend/internal/storage"
)

type Handler struct {
	configurationService configuration.Service
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(configurationService configuration.Service) http.Handler {
	handler := Handler{configurationService: configurationService}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handler.handleHealth)
	mux.HandleFunc("GET /api/configuration/current", handler.handleGetCurrentConfiguration)
	mux.HandleFunc("PUT /api/configuration/current", handler.handlePutCurrentConfiguration)
	mux.HandleFunc("POST /api/configuration/generate", handler.handleGenerateConfiguration)

	return mux
}

func (h Handler) handleHealth(responseWriter http.ResponseWriter, _ *http.Request) {
	writeJSON(responseWriter, http.StatusOK, map[string]string{"status": "ok"})
}

func (h Handler) handleGetCurrentConfiguration(responseWriter http.ResponseWriter, request *http.Request) {
	current, err := h.configurationService.GetCurrent(request.Context())
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

func (h Handler) handlePutCurrentConfiguration(responseWriter http.ResponseWriter, request *http.Request) {
	var current configuration.ActiveConfiguration
	if err := decodeJSON(request, &current); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.configurationService.SaveCurrent(request.Context(), current); err != nil {
		if configuration.IsValidationError(err) {
			writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		writeJSON(responseWriter, http.StatusInternalServerError, errorResponse{Error: "failed to save active configuration"})
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}

func (h Handler) handleGenerateConfiguration(responseWriter http.ResponseWriter, request *http.Request) {
	var generateRequest configuration.GenerateRequest
	if err := decodeJSON(request, &generateRequest); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	configurationResult, err := h.configurationService.Generate(request.Context(), generateRequest)
	if err != nil {
		switch {
		case configuration.IsValidationError(err):
			writeJSON(responseWriter, http.StatusBadRequest, errorResponse{Error: err.Error()})
		case errors.Is(err, configuration.ErrGenerateNotImplemented):
			writeJSON(responseWriter, http.StatusNotImplemented, errorResponse{Error: err.Error()})
		default:
			writeJSON(responseWriter, http.StatusInternalServerError, errorResponse{Error: "failed to generate configuration"})
		}

		return
	}

	writeJSON(responseWriter, http.StatusOK, configurationResult)
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
