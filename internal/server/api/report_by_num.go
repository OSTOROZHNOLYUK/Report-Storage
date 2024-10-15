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

// ReportGetter - интерфейс для доступа к заявке по её номеру.
type ReportGetter interface {
	ReportByNum(context.Context, int) (storage.Report, error)
}

// ReportByNum обрабатывает запрос на получение заявки по её номеру.
func ReportByNum(l *slog.Logger, st ReportGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportByNum"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to receive report by number")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса.
		num, err := number(r)
		if err != nil {
			log.Error("invalid report number", logger.Err(err))
			http.Error(w, "invalid report number", http.StatusBadRequest)
			return
		}

		// Запрос в базу данных.
		report, err := st.ReportByNum(r.Context(), num)
		if err != nil {
			log.Error("cannot find report", logger.Err(err))
			if errors.Is(err, storage.ErrReportNotFound) {
				http.Error(w, "report not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		err = json.NewEncoder(w).Encode(report)
		if err != nil {
			log.Error("cannot encode report", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("report sent successfully")
	}
}
