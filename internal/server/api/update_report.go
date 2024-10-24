package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/notifications"
	"Report-Storage/internal/reports"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// ReportUpdater - интерфейс для обновления всех полей заявки.
type ReportUpdater interface {
	UpdateReport(ctx context.Context, rep storage.Report) (storage.Report, error)
}

// UpdateReport обрабатывает запрос на обновление заявки по
// уникальному номеру.
func UpdateReport(l *slog.Logger, st ReportUpdater, s3 reports.FileSaver, notify *notifications.SMTP) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.UpdateReport"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to update report data")

		// Установка типа контента для ответа.
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

		// Запрос в базу данных.
		// Метод UpdateReport возвращает заявку ДО ее изменения. Это необходимо
		// для сравнения некоторых полей и удаления неиспользуемых файлов.
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

		// Проверка медиа файлов.
		// Если в измененной заявке меньше медиа файлов, чем до изменения,
		// то находим разницу и удаляем неиспользуемые файлы из S3 хранилища.
		if len(origin.Media) > len(report.Media) {
			diff := reports.SliceDiff(origin.Media, report.Media)
			if len(diff) > 0 {
				log.Debug("removing media files")
				go reports.RemoveFiles(log, diff, s3)
			}
		}

		// Проверка изменения статуса и отправка уведомления об этом.
		if origin.Status != report.Status && report.Contacts.Email != "" {
			go func() {
				err := notifications.StatusChanged(notify, report.Contacts.Email, statusString(report.Status))
				if err != nil {
					log.Error("failed to send notification to email", logger.Err(err))
				}
			}()
		}

		// Кодирование ответа в JSON.
		err = json.NewEncoder(w).Encode(report)
		if err != nil {
			log.Error("cannot encode report", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("report sent successfully")
	}
}
