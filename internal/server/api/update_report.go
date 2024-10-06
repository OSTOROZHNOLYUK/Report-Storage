package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/reports"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// ReportUpdater - интерфейс для обновления всех полей заявки.
type ReportUpdater interface {
	UpdateReport(ctx context.Context, rep storage.Report) (storage.Report, error)
}

// UpdateReport обрабатывает запрос на обновление заявки по уникальному номеру.
func UpdateReport(l *slog.Logger, st ReportUpdater, s3 reports.FileSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.UpdateReport"
		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Установка типа контента для ответа
		w.Header().Set("Content-Type", "application/json")

		// Декодируем тело запроса в структуру.
		var report storage.Report
		if err := render.DecodeJSON(r.Body, &report); err != nil {
			log.Error("failed to decode JSON", logger.Err(err))
			http.Error(w, "invalid report data", http.StatusBadRequest)
			return
		}

		// Валидируем поля запроса.
		valid := validator.New()
		err := valid.Struct(report)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("validation failed", logger.Err(validateErr))
			http.Error(w, "invalid report data", http.StatusBadRequest)
			return
		}
		report.Updated = time.Now()
		report.Geo.Type = "Point"
		log.Debug("json input decoded and validated successfully")

		// Метод UpdateReport возвращает заявку ДО ее изменения. Это необходимо
		// для сравнения полей Media и удаления неиспользуемых файлов.
		origin, err := st.UpdateReport(r.Context(), report)
		if err != nil {
			log.Error("failed to update report", logger.Err(err))
			if errors.Is(err, storage.ErrReportNotFound) {
				http.Error(w, "report not found", http.StatusNotFound)
				return
			}
			if errors.Is(err, storage.ErrIncorrectID) {
				http.Error(w, "invalid report data", http.StatusBadRequest)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if !reflect.DeepEqual(report.Media, origin.Media) {

			// TODO: написать алгортм удаления неиспользуемых файлов из объектного
			//       хранилища.
			log.Debug("removing files from S3")
		}

		err = json.NewEncoder(w).Encode(report)
		if err != nil {
			log.Error("cannot encode report", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Debug("report sent successfully")
	}
}
