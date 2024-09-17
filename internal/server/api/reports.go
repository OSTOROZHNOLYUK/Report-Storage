package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Reporter interface {
	Reports(ctx context.Context, status []storage.Status) ([]storage.Report, error)
}

func Reports(log *slog.Logger, st Reporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.Reports"

		log.Info("new request to receive all reports")

		// TODO: получение статусов.

		var status []storage.Status
		ctx := r.Context()
		reports, err := st.Reports(ctx, status)
		if err != nil {
			log.Error("cannot receive all reports from DB", logger.Err(err), slog.String("op", operation))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("all reports received from DB")

		err = json.NewEncoder(w).Encode(reports)
		if err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err), slog.String("op", operation))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("all reports encoded")
		w.Header().Set("Content-Type", "application/json")
	}
}
