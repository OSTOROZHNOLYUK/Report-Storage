package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/reports"
	"Report-Storage/internal/storage"
	"context"

	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ReportUpdater interface {
	UpdateReport(ctx context.Context, rep storage.Report) (storage.Report, error)
}

// UpdateReport обрабатывает запрос на обновление заявки по уникальному номеру.
func UpdateReport(l *slog.Logger, st ReportUpdater, s3 reports.FileSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.UpdateReport"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		numStr := chi.URLParam(r, "num")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.Error("invalid report number", logger.Err(err))
			http.Error(w, "invalid report number", http.StatusBadRequest)
			return
		}

		var report storage.Report
		if err := render.Bind(r, &report); err != nil {
			log.Error("failed to bind JSON", logger.Err(err))
			http.Error(w, "failed to bind JSON", http.StatusBadRequest)
			return
		}

		report.Number = int64(num)
		origin, err := st.UpdateReport(r.Context(), report)
		if err != nil {
			if errors.Is(err, storage.ErrReportNotFound) {
				http.Error(w, "report not found", http.StatusNotFound)
			} else {
				log.Error("failed to update report", logger.Err(err))
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}

		render.JSON(w, r, origin)
	}
}
