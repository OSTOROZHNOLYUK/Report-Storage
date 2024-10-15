package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"errors"
	"log/slog"
	"net/http"
)

// ReportRemover - интерфейс для удаления заявки.
type ReportRemover interface {
	DeleteByNum(ctx context.Context, num int) error
}

// DeleteReport обрабатывает запрос на удаление заявки по её номеру.
func DeleteReport(l *slog.Logger, st ReportRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.DeleteReport"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to delete report")

		// Получение параметров запроса.
		num, err := number(r)
		if err != nil {
			log.Error("invalid report number", logger.Err(err))
			http.Error(w, "invalid report number", http.StatusBadRequest)
			return
		}

		// Запрос в базу данных.
		err = st.DeleteByNum(r.Context(), num)
		if err != nil {
			log.Error("cannot delete report", logger.Err(err))
			if errors.Is(err, storage.ErrReportNotFound) {
				http.Error(w, "report not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Запись кода ответа.
		w.WriteHeader(http.StatusNoContent)
		log.Debug("report deleted successfully")
	}
}
