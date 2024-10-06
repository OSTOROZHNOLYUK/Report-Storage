package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

// RejectRemover - интерфейс для удаления отклоненных заявок.
type RejectRemover interface {
	DeleteRejected(ctx context.Context) (int, error)
}

// DeleteRejected обрабатывает запрос на удаление всех отклоненных заявок.
func DeleteRejected(l *slog.Logger, st RejectRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.DeleteRejected"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to delete rejected reports")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		ctx := r.Context()
		count, err := st.DeleteRejected(ctx)
		if err != nil {
			log.Error("cannot delete rejected reports", logger.Err(err))
			if errors.Is(err, storage.ErrReportNotFound) {
				http.Error(w, "no rejected reports found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("rejected reports deleted succesfully", slog.Int("count", count))

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Deleted rejected reports: " + strconv.Itoa(count)))
		if err != nil {
			log.Error("cannot write response", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
}
