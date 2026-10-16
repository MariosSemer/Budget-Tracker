package main

import (
	"net/http"

	"BudgetTrucker/internal/httpapi"
	slogger "BudgetTrucker/internal/logger"
)

func main() {

	log := slogger.New()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", httpapi.HealthHandler)

	log.Info("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Error("server crashed", "error", err)
	}
}
