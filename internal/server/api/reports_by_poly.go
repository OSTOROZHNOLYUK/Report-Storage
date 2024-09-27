package api

import (
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// ReportsByQuadHandler - интерфейс для работы с хранилищем.
type ReportsByPolyHandler interface {
	ReportsByPoly(ctx context.Context, poly [][2]float64, status []storage.Status) ([]storage.Report, error)
}

// ReportsByQuad обрабатывает запросы для получения заявок в границах многоугольника.
func ReportsByPoly(log *slog.Logger, st ReportsByPolyHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Quad   [][2]float64     `json:"quad"`
			Status []storage.Status `json:"status,omitempty"`
		}

		// Декодируем JSON из тела запроса.
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// Получаем заявки из хранилища.
		reports, err := st.ReportsByPoly(r.Context(), input.Quad, input.Status)
		if err != nil {
			// Обработка ошибок.
			if err == storage.ErrArrayNotFound {
				http.Error(w, "no reports found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Кодируем результаты в JSON и отправляем ответ.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(reports); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
