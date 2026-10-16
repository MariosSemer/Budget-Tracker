package httpapi

import (
	"fmt"
	"net/http"
	//"BudgetTrucker/internal/slogger"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}
