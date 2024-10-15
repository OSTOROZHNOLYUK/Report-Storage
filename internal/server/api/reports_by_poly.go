package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// input - структура тела запроса с координатами вершин
// многоугольника.
type polygon struct {
	Quad [][2]float64 `json:"quad"`
}

// ReportsByPolyInterface - интерфейс для получения заявок
// в границах многоугольника.
type ReportsByPolyInterface interface {
	ReportsByPoly(ctx context.Context, poly [][2]float64, status []storage.Status) ([]storage.Report, error)
}

// ReportsByPoly обрабатывает запросы для получения заявок
// в границах многоугольника.
func ReportsByPoly(l *slog.Logger, st ReportsByPolyInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportsByPoly"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to receive reports by polygon")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Декодирование JSON из тела запроса.
		var input polygon
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Error("cannot decode json to polygon struct", logger.Err(err))
			http.Error(w, "invalid request JSON", http.StatusBadRequest)
			return
		}

		// Проверка количества вершин многоугольника.
		if len(input.Quad) < 3 {
			log.Error("vertex count less than 3")
			http.Error(w, "invalid request JSON", http.StatusBadRequest)
			return
		}

		// Получение статусов.
		statusParam := r.URL.Query().Get("status")
		status := splitStatus(statusParam)

		// Запрос в базу данных.
		reports, err := st.ReportsByPoly(r.Context(), input.Quad, status)
		if err != nil {
			log.Error("failed to get reports by polygon", logger.Err(err))
			if errors.Is(err, storage.ErrArrayNotFound) {
				http.Error(w, "no reports found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		if err := json.NewEncoder(w).Encode(reports); err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("reports by polygon encoded and sent successfully")
	}
}
