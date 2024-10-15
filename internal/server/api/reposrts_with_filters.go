package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// ReportsFilterer - интерфейс для получения заявок с фильтром.
type ReportsFilterer interface {
	ReportsWithFilter(ctx context.Context, fl storage.Filter) ([]storage.Report, error)
}

// ReportsWithFilters обрабатывает запрос на получение N последних
// заявок с фильтрами.
func ReportsWithFilters(l *slog.Logger, st ReportsFilterer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportsWithFilters"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to receive filtered report list")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса.
		s := r.URL.Query().Get("status")
		status := splitStatus(s)
		filter := storage.Filter{
			Count:  count(r),
			Sort:   sort(r),
			Status: status,
		}

		// Запрос в базу данных.
		reports, err := st.ReportsWithFilter(r.Context(), filter)
		if err != nil {
			log.Error("failed to get reports with filter", logger.Err(err))
			if errors.Is(err, storage.ErrArrayNotFound) {
				http.Error(w, "no reports found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		err = json.NewEncoder(w).Encode(reports)
		if err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("filtered reports encoded and sent successfully")
	}
}
