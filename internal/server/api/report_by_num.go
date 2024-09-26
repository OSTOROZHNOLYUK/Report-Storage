package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
		)
		log.Info("request to get report by num")

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
			if err == storage.ErrReportNotFound {
				http.Error(w, "report not found", http.StatusNotFound)
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(report)
		if err != nil {
			log.Error("cannot encode report", logger.Err(err))
			http.Error(w, "failed to encode report", http.StatusInternalServerError)
			return
		}

		log.Debug("report sent successfully")
	}
}
