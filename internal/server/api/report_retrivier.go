package api

import (
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ReportRetriever - интерфейс для получения заявки по ObjectID
type ReportRetriever interface {
	ReportByID(ctx context.Context, id string) (storage.Report, error)
}

// GetReportByID обрабатывает запрос на получение заявки по ObjectID.
func GetReportByID(log *slog.Logger, st ReportRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		report, err := st.ReportByID(r.Context(), id)
		if err != nil {
			if err.Error() == "report not found" {
				http.Error(w, "report not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Error("database error:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	}
}
