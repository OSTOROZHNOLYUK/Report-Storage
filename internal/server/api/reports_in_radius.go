package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"Report-Storage/internal/storage"
)

//type ReportsInRadiusHandler interface {
//	ReportsInRadius(ctx context.Context, x, y float64, r int, statuses []storage.Status) ([]storage.Report, error)
//}

type ReportsByRadiusHandler interface {
	GetReportsByRadius(ctx context.Context, r int, p storage.Geo, status []storage.Status) ([]storage.Report, error)
}

func ReportsByRadius(log *slog.Logger, st ReportsByRadiusHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.ReportsInRadius"
		log := log.With(slog.String("op", operation))
		log.Info("Начало обработки запроса")

		// Получение параметров запроса
		xParam := r.URL.Query().Get("x")
		yParam := r.URL.Query().Get("y")
		radiusParam := r.URL.Query().Get("radius")
		statusesParam := r.URL.Query().Get("statuses")

		// Проверка на отсутствие параметров
		if xParam == "" || yParam == "" || radiusParam == "" {
			http.Error(w, "Недостаточно параметров", http.StatusBadRequest)
			log.Error("Недостаточно параметров")
			return
		}

		// Преобразование координат и радиуса в нужный формат
		x, err := strconv.ParseFloat(xParam, 64)
		if err != nil {
			http.Error(w, "Неверный формат параметра x", http.StatusBadRequest)
			log.Error("Неверный формат параметра x", slog.String("error", err.Error()))
			return
		}

		y, err := strconv.ParseFloat(yParam, 64)
		if err != nil {
			http.Error(w, "Неверный формат параметра y", http.StatusBadRequest)
			log.Error("Неверный формат параметра y", slog.String("error", err.Error()))
			return
		}

		radius, err := strconv.Atoi(radiusParam)
		if err != nil {
			http.Error(w, "Неверный формат параметра radius", http.StatusBadRequest)
			log.Error("Неверный формат параметра radius", slog.String("error", err.Error()))
			return
		}

		// Преобразование статусов в массив
		statuses, err := parseStatuses(statusesParam)
		if err != nil {
			// обработка ошибки
			fmt.Println("Ошибка преобразования статусов:", err)
			return
		}
		// Вызов метода для получения заявок

		position := storage.Geo{Type: "Point",
			Coordinates: [2]float64{x, y}}
		reports, err := st.GetReportsByRadius(r.Context(), radius, position, statuses)
		if err != nil {
			http.Error(w, "Ошибка при получении заявок: "+err.Error(), http.StatusInternalServerError)
			log.Error("Ошибка при получении заявок", slog.String("error", err.Error()))
			return
		}

		// Кодирование ответа
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reports); err != nil {
			http.Error(w, "Ошибка при кодировании ответа: "+err.Error(), http.StatusInternalServerError)
			log.Error("Ошибка при кодировании ответа", slog.String("error", err.Error()))
			return
		}

		log.Info("Обработка запроса завершена успешно")
	}
}

// parseStatuses преобразует строку статусов в срез статусов
func parseStatuses(statusesParam string) ([]storage.Status, error) {
	var statuses []storage.Status
	if statusesParam == "" {
		return statuses, nil
	}

	statusNames := strings.Split(statusesParam, ",")
	for _, statusName := range statusNames {
		status, err := storage.StatusFromString(strings.TrimSpace(statusName))
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}
