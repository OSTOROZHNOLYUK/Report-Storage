package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

// StatisticRetriever - интерфейс для получения статистики заявок.
type StatisticRetriever interface {
	Statistic(ctx context.Context) (storage.Statistic, error)
}

// Statistic обрабатывает запрос на получение статистики по всем заявкам.
func Statistic(l *slog.Logger, st StatisticRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.GetStatistic"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to receive statistics")

		ctx := r.Context()
		stats, err := st.Statistic(ctx)
		if err != nil {
			log.Error("cannot retrieve statistics", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Error("cannot encode statistics", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Debug("statistics sent successfully")
	}
}
