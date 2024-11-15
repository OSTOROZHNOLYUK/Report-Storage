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

// Reporter - интерфейс для БД в обработчике Reports.
type Reporter interface {
	Reports(ctx context.Context, status []storage.Status) ([]storage.Report, error)
}

// Reports обрабатывает запрос на получение всех заявок с возможностью
// фильтрации по статусам. Статусы принимаются query параметром status
// со значениями с виде чисел через запятую. Числа соответствуют
// константам из пакета storage.
func Reports(l *slog.Logger, st Reporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.Reports"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to receive all reports with status")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса.
		s := r.URL.Query().Get("status")
		status := splitStatus(s)

		// Запрос в базу данных.
		reports, err := st.Reports(r.Context(), status)
		if err != nil {
			log.Error("cannot receive all reports", logger.Err(err))
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
		log.Debug("all reports encoded and sent successfully")
	}
}
