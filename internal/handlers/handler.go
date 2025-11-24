package handlers

import (
	"encoding/json"
	"mini-quicko/internal/service"
	"mini-quicko/pkg/models"
	"net/http"
	"time"
)

type Handler struct {
	analyzer *service.Analyzer
}

func NewHandler(analyzer *service.Analyzer) *Handler {
	return &Handler{
		analyzer: analyzer,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) AnalyzeProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed", "Only GET requests are allowed")
		return
	}

	productID := r.URL.Query().Get("id")
	if productID == "" {
		respondError(w, http.StatusBadRequest, "missing_product_id", "Product ID is required")
		return
	}

	analysis, err := h.analyzer.AnalyzeProduct(productID)
	if err != nil {
		respondError(w, http.StatusNotFound, "analysis_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, analysis)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed", "Only GET requests are allowed")
		return
	}

	productID := r.URL.Query().Get("id")
	if productID == "" {
		respondError(w, http.StatusBadRequest, "missing_product_id", "Product ID is required")
		return
	}

	history, err := h.analyzer.GetHistory(productID)
	if err != nil {
		respondError(w, http.StatusNotFound, "history_not_found", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, history)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, errorCode, message string) {
	response := models.ErrorResponse{
		Error:   errorCode,
		Message: message,
	}
	respondJSON(w, status, response)
}
