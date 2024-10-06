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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// // Структура для обновления статуса
// type updateStatusRequest struct {
// 	NewStatus storage.Status `json:"new"`
// }

// ReportStatusUpdater - интерфейс для обновления статуса заявки.
type ReportStatusUpdater interface {
	UpdateStatus(ctx context.Context, num int, status storage.Status) (storage.Report, error)
}

// UpdateStatusReport обрабатывает запрос для изменения статуса заявки.
func UpdateStatusReport(l *slog.Logger, st ReportStatusUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.UpdateStatusReport"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to update report status")

		// Установка типа контента для ответа
		w.Header().Set("Content-Type", "application/json")

		numStr := chi.URLParam(r, "num")
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 1 {
			log.Error("invalid report number", logger.Err(err))
			http.Error(w, "invalid report number", http.StatusBadRequest)
			return
		}

		// Получаем значение нового статуса.
		s := r.URL.Query().Get("new")
		if s == "" {
			log.Error("empty new status")
			http.Error(w, "incorrect new status", http.StatusBadRequest)
			return
		}
		status := splitStatus(s)[0]
		if status < 1 || status > 5 {
			log.Error("incorrect new status", slog.Int("status", int(status)))
			http.Error(w, "incorrect new status", http.StatusBadRequest)
			return
		}

		// var req updateStatusRequest
		// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// 	log.Error("invalid request body", logger.Err(err))
		// 	http.Error(w, "invalid request body", http.StatusBadRequest)
		// 	return
		// }

		// Обновляем статус заявки
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

		err = json.NewEncoder(w).Encode(report)
		if err != nil {
			log.Error("cannot encode report", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Debug("report status updated successfully")
	}
}
