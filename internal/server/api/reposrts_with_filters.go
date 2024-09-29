package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

// ReportsFilterer - интерфейс для получения заявок с фильтром.
type ReportsFilterer interface {
	ReportsWithFilter(ctx context.Context, fl storage.Filter) ([]storage.Report, error)
}

// ReportsWithFilters обрабатывает запрос на получение N последних заявок с фильтрами.
func ReportsWithFilters(l *slog.Logger, st ReportsFilterer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportsWithFilters"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to receive filtered report list")

		// Установка типа контента для ответа
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров `n`, `status` и `sort` из запроса
		nStr := r.URL.Query().Get("n")
		statusStr := r.URL.Query().Get("status")
		sortStr := r.URL.Query().Get("sort")

		// Разбор параметра n
		n := 20 // Значение по умолчанию
		if nStr != "" {
			var err error
			n, err = strconv.Atoi(nStr)
			if err != nil || n < 1 {
				log.Warn("Invalid parameter n, using default value", slog.Int("default_n", n))
				// Если параметр некорректен, используем значение по умолчанию
				n = 20
			}
		}

		// Разбор параметра sort
		sort := -1 // Значение по умолчанию
		if sortStr != "" {
			var err error
			sort, err = strconv.Atoi(sortStr)
			if err != nil || (sort != 1 && sort != -1) {
				log.Warn("Invalid parameter sort, using default value", slog.Int("default_sort", sort))
				// Если параметр некорректен, используем значение по умолчанию
				sort = -1
			}
		}

		// Разбор параметра status
		status := splitStatus(statusStr)

		// Фильтр для запроса
		filter := storage.Filter{
			Count:  n,
			Sort:   sort,
			Status: status,
		}

		ctx := r.Context()
		reports, err := st.ReportsWithFilter(ctx, filter)
		if err != nil {
			log.Error("failed to get reports with filter", logger.Err(err))
			if errors.Is(err, storage.ErrArrayNotFound) {
				http.Error(w, "no reports found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(reports)
		if err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("filtered reports encoded and sent successfully")
	}
}
