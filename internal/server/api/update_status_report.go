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

// ReportStatusUpdater - интерфейс для обновления статуса заявки.
type ReportStatusUpdater interface {
	UpdateStatus(ctx context.Context, num int, status storage.Status) (storage.Report, error)
}

// UpdateStatusReport обрабатывает запрос для изменения статуса заявки.
func UpdateStatusReport(l *slog.Logger, st ReportStatusUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.UpdateStatusReport"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to update report status")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса.
		num, err := number(r)
		if err != nil {
			log.Error("invalid report number", logger.Err(err))
			http.Error(w, "invalid report number", http.StatusBadRequest)
			return
		}
		status, err := newStatus(r)
		if err != nil {
			log.Error("incorrect new status", logger.Err(err))
			http.Error(w, "incorrect new statusr", http.StatusBadRequest)
			return
		}

		// Запрос в базу данных.
		report, err := st.UpdateStatus(r.Context(), num, status)
		if err != nil {
			log.Error("cannot update report status", logger.Err(err))
			if errors.Is(err, storage.ErrReportNotFound) {
				http.Error(w, "report not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// TODO: вставить нотификацию по контактам.

		// Кодирование ответа в JSON.
		err = json.NewEncoder(w).Encode(report)
		if err != nil {
			log.Error("cannot encode report", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("report status updated successfully")
	}
}
