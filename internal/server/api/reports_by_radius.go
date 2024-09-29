package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
)

// ReportsByRadiusHandler - интерфейс для получения заявок в радиусе от точки.
type ReportsByRadiusHandler interface {
	ReportsByRadius(ctx context.Context, r int, p storage.Geo, status []storage.Status) ([]storage.Report, error)
}

// ReportsByRadius обрабатывает запрос на получение заявок в радиусе от точки.
func ReportsByRadius(l *slog.Logger, st ReportsByRadiusHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportsInRadius"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to receive reports by radius")

		// Установка типа контента для ответа
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса
		xParam := r.URL.Query().Get("x")
		yParam := r.URL.Query().Get("y")
		radiusParam := r.URL.Query().Get("r")
		statusesParam := r.URL.Query().Get("status")

		// Проверка на отсутствие параметров
		if xParam == "" || yParam == "" || radiusParam == "" {
			http.Error(w, "invalid parameters", http.StatusBadRequest)
			log.Error("insufficient parameters")
			return
		}

		// Преобразование координат и радиуса в нужный формат
		x, err := strconv.ParseFloat(xParam, 64)
		if err != nil {
			http.Error(w, "invalid parameters", http.StatusBadRequest)
			log.Error("invalid parameter x format", logger.Err(err))
			return
		}

		y, err := strconv.ParseFloat(yParam, 64)
		if err != nil {
			http.Error(w, "invalid parameters", http.StatusBadRequest)
			log.Error("invalid parameter y format", logger.Err(err))
			return
		}

		radius, err := strconv.Atoi(radiusParam)
		if err != nil {
			http.Error(w, "invalid parameters", http.StatusBadRequest)
			log.Error("Invalid parameter radius format", logger.Err(err))
			return
		}

		// Разбор параметра status
		status := splitStatus(statusesParam)
		position := storage.Geo{
			Type:        "Point",
			Coordinates: [2]float64{x, y},
		}

		// // Преобразование статусов в массив
		// statuses, err := parseStatuses(statusesParam)
		// if err != nil {
		// 	// обработка ошибки
		// 	fmt.Println("Status conversion error:", err)
		// 	return
		// }
		// Вызов метода для получения заявок

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

		// Кодирование ответа
		err = json.NewEncoder(w).Encode(reports)
		if err != nil {
			log.Error("cannot encode reports to ResponseWriter", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("reports by radius encoded and sent successfully")
	}
}

// parseStatuses преобразует строку статусов в срез статусов
// func parseStatuses(statusesParam string) ([]storage.Status, error) {
// 	var statuses []storage.Status
// 	if statusesParam == "" {
// 		return statuses, nil
// 	}

// 	statusNames := strings.Split(statusesParam, ",")
// 	for _, statusName := range statusNames {
// 		status, err := storage.StatusFromString(strings.TrimSpace(statusName))
// 		if err != nil {
// 			return nil, err
// 		}
// 		statuses = append(statuses, status)
// 	}
// 	return statuses, nil
// }
