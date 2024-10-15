package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// RejectRemover - интерфейс для удаления отклоненных заявок.
type RejectRemover interface {
	DeleteRejected(ctx context.Context) (int, error)
}

// DeleteRejected обрабатывает запрос на удаление всех отклоненных заявок.
func DeleteRejected(l *slog.Logger, st RejectRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.DeleteRejected"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to delete rejected reports")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Запрос в базу данных.
		count, err := st.DeleteRejected(r.Context())
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

		// Запись ответа в text/plain.
		str := fmt.Sprintf("Deleted rejected reports: %d", count)
		_, err = w.Write([]byte(str))
		if err != nil {
			log.Error("cannot write response", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("response sent successfully")
	}
}
