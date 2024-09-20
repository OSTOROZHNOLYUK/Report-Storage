package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
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

		log := l.With(
			slog.String("op", operation),
		)
		log.Info("request to receive all reports with status")

		w.Header().Set("Content-Type", "application/json")

		s := r.URL.Query().Get("status")
		status := splitStatus(s)

		ctx := r.Context()
		reports, err := st.Reports(ctx, status)
		if err != nil {
			log.Error("cannot receive all reports from DB", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("all reports received from DB")

		err = json.NewEncoder(w).Encode(reports)
		if err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("all reports encoded")
		w.WriteHeader(http.StatusOK)
	}
}
