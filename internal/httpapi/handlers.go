package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

type AddHandler struct {
	logger *slog.Logger
}

func (h *HealthHandler) AddHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprintln(w, "ok"); err != nil {
		h.Logger.Error("Health check write failed", "error", err)
	}
}
