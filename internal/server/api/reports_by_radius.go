package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
)

// ReportsByRadiusInterface - интерфейс для получения заявок
// в радиусе от точки.
type ReportsByRadiusInterface interface {
	ReportsByRadius(ctx context.Context, r int, p storage.Geo, status []storage.Status) ([]storage.Report, error)
}

// ReportsByRadius обрабатывает запрос на получение заявок
// в радиусе от точки.
func ReportsByRadius(l *slog.Logger, st ReportsByRadiusInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportsInRadius"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to receive reports by radius")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса.
		position, radius, err := point(r)
		if err != nil {
			log.Error("failed to get correct parameters", logger.Err(err))
			http.Error(w, "invalid parameters", http.StatusBadRequest)
			return
		}
		s := r.URL.Query().Get("status")
		status := splitStatus(s)

		// Запрос в базу данных.
		reports, err := st.ReportsByRadius(r.Context(), radius, position, status)
		if err != nil {
			log.Error("failed to get reports by radius", logger.Err(err))
			if errors.Is(err, storage.ErrArrayNotFound) {
				http.Error(w, "no reports found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		err = json.NewEncoder(w).Encode(reports)
		if err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("reports by radius encoded and sent successfully")
	}
}
