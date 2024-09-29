package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ReportRetriever - интерфейс для получения заявки по ObjectID
type ReportRetriever interface {
	ReportByID(ctx context.Context, id string) (storage.Report, error)
}

// ReportByID обрабатывает запрос на получение заявки по ObjectID.
func ReportByID(l *slog.Logger, st ReportRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportByID"

		log := l.With(slog.String("op", operation))
		log.Info("request to receive report by objectid")

		// Установка типа контента для ответа
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParam(r, "id")
		fmt.Println(id)
		if id == "" {
			log.Error("empty id value")
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		report, err := st.ReportByID(r.Context(), id)
		if err != nil {
			log.Error("cannot find report", logger.Err(err))
			if errors.Is(err, storage.ErrIncorrectID) {
				http.Error(w, "invalid objectid", http.StatusBadRequest)
				return
			}
			if errors.Is(err, storage.ErrReportNotFound) {
				http.Error(w, "report not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
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
