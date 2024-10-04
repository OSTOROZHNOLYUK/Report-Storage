package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Структура для обновления статуса
type updateStatusRequest struct {
	NewStatus storage.Status `json:"new"`
}

// UpdateStatusReport обрабатывает запрос для изменения статуса заявки.
func UpdateStatusReport(l *slog.Logger, st ReportGetter) http.HandlerFunc {
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

		var req updateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("invalid request body", logger.Err(err))
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Обновляем статус заявки
		ctx := r.Context()
		report, err := st.ReportByNum(ctx, num)
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
