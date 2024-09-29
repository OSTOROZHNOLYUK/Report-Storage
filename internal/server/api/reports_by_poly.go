package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// ReportsByPolyHandler - интерфейс для получения заявок в границах многоугольника.
type ReportsByPolyHandler interface {
	ReportsByPoly(ctx context.Context, poly [][2]float64, status []storage.Status) ([]storage.Report, error)
}

// ReportsByPoly обрабатывает запросы для получения заявок в границах многоугольника.
func ReportsByPoly(l *slog.Logger, st ReportsByPolyHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportsByPoly"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to receive reports by polygon")

		// Установка типа контента для ответа
		w.Header().Set("Content-Type", "application/json")

		// input - структура запроса.
		var input struct {
			Quad [][2]float64 `json:"quad"`
			// Status []storage.Status `json:"status,omitempty"`
		}

		// Декодируем JSON из тела запроса.
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Error("cannot decode json to input struct", logger.Err(err))
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// Получение статусов.
		statusParam := r.URL.Query().Get("status")
		status := splitStatus(statusParam)

		// Получаем заявки из хранилища.
		reports, err := st.ReportsByPoly(r.Context(), input.Quad, status)
		if err != nil {
			// Обработка ошибок.
			log.Error("failed to get reports by polygon", logger.Err(err))
			if errors.Is(err, storage.ErrArrayNotFound) {
				http.Error(w, "no reports found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодируем результаты в JSON и отправляем ответ.
		if err := json.NewEncoder(w).Encode(reports); err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
		log.Debug("reports by polygon encoded and sent successfully")
	}
}
