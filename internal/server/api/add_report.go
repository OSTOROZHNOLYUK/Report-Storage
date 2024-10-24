package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/notifications"
	"Report-Storage/internal/reports"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

const (
	// maxMemory - максимальный размер тела запроса.
	maxMemory int64 = 30 << 20
)

// AddReport обрабатывает запрос на добавление новой заявки в хранилище.
// При успехе возвращает код 201 и уникальный номер заявки.
func AddReport(l *slog.Logger, st reports.ReportAdder, s3 reports.FileSaver, notify *notifications.SMTP) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.UploadFiles"

		// Настройка логирования.
		log := logger.Handler(l, operation, r)
		log.Info("request to add new report")

		// Проверка заголовка Content-Type на значение multipart/form-data.
		ct := strings.ToLower(r.Header.Get("Content-Type"))
		if !strings.Contains(ct, "multipart/form-data") {
			log.Error("content-type is not multipart/form-data", slog.String("Content-Type", ct))
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		// Суммарный размер всех загружаемых файлов не более 30 Мб.
		r.Body = http.MaxBytesReader(w, r.Body, maxMemory)
		r.ParseMultipartForm(maxMemory + 512)
		defer func() {
			err := r.MultipartForm.RemoveAll()
			if err != nil {
				log.Error("cannot remove temporary multipart form files", logger.Err(err))
			}
		}()

		// Проверка наличия файлов и строковых значений.
		if len(r.MultipartForm.Value) == 0 {
			log.Error("json not found in the request body")
			http.Error(w, "incorrect report data", http.StatusBadRequest)
			return
		}
		if len(r.MultipartForm.File) == 0 || len(r.MultipartForm.File) > 5 {
			log.Error("incorrect files count")
			http.Error(w, "incorrect report data", http.StatusBadRequest)
			return
		}

		// Получение сформированной структуры заявки и кода. Если code
		// не равно 200, то возвращаем ошибку.
		report, code := reports.Build(l, s3, r)
		switch code {
		case http.StatusBadRequest:
			http.Error(w, "incorrect report data", http.StatusBadRequest)
			return
		case http.StatusInternalServerError:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		case http.StatusUnsupportedMediaType:
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		// TODO: добавить пул для конвертации фото.
		//
		// Принудительный возврат аллоцированной памяти системе.
		defer debug.FreeOSMemory()

		log.Debug("request body parsed succefully")

		// Получение контекста запроса.
		ctx := r.Context()

		// Получение нового номера заявки и запись его в структуру заявки.
		// ObjectID будет сгенерирован в методе БД AddReport. В случае
		// ошибки удаляем загруженные файлы из S3 хранилища.
		newNum, err := st.CounterInc(ctx)
		if err != nil {
			go reports.RemoveFiles(l, report.Media, s3)
			log.Error("cannot receive new ID", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		report.Number = int64(newNum)

		// Добавление сформированной заявки в БД. В случае ошибки удаляем
		// загруженные файлы из S3 хранилища.
		err = st.AddReport(ctx, report)
		if err != nil {
			go reports.RemoveFiles(l, report.Media, s3)
			log.Error("cannot add report to DB", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("new report added successfully")

		// Отправка уведомления о создании новой заявки.
		if report.Contacts.Email != "" {
			go func() {
				err := notifications.NewReport(notify, report.Contacts.Email)
				if err != nil {
					log.Error("failed to send notification to email", logger.Err(err))
				}
			}()
		}

		// Запись ответа в text/plain и установка кода 201.
		render.Status(r, http.StatusCreated)
		render.PlainText(w, r, strconv.Itoa(int(report.Number)))
		log.Debug("new report number sent successfully")
	}
}
