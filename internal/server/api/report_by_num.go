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

// ReportGetter - интерфейс для доступа к заявке по её номеру.
type ReportGetter interface {
	ReportByNum(context.Context, int) (storage.Report, error)
}

// ReportByNum обрабатывает запрос на получение заявки по её номеру.
func ReportByNum(l *slog.Logger, st ReportGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportByNum"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to receive report by num")

		// Установка типа контента для ответа
		w.Header().Set("Content-Type", "application/json")

		numStr := chi.URLParam(r, "num")
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 1 {
			log.Error("invalid report number", logger.Err(err))
			http.Error(w, "invalid report number", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		report, err := st.ReportByNum(ctx, num)
		if err != nil {
			log.Error("cannot find report", logger.Err(err))
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

		log.Debug("report sent successfully")
	}
}
