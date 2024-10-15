package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"net/http"

	"log/slog"
)

// StatisticRetriever - интерфейс для получения статистики заявок.
type StatisticRetriever interface {
	Statistic(ctx context.Context) (storage.Statistic, error)
}

// Statistic обрабатывает запрос на получение статистики по всем заявкам.
func Statistic(l *slog.Logger, st StatisticRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.GetStatistic"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to receive statistics")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Запрос в базу данных.
		stats, err := st.Statistic(r.Context())
		if err != nil {
			log.Error("cannot retrieve statistics", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Error("cannot encode statistics", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("statistics sent successfully")
	}
}
