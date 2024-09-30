package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/reports"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	// maxMemory - максимальный размер тела запроса.
	maxMemory int64 = 30 << 20
)

// AddReport обрабатывает запрос на добавление новой заявки в хранилище.
// При успехе возвращает код 201 и уникальный номер заявки.
func AddReport(l *slog.Logger, st reports.ReportAdder, s3 reports.FileSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.UploadFiles"

		log := l.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("request to add new report")

		// Проверяем заголовок Content-Type на значение multipart/form-data.
		ct := strings.ToLower(r.Header.Get("Content-Type"))
		if !strings.Contains(ct, "multipart/form-data") {
			log.Error("content-type is not multipart/form-data", slog.String("Content-Type", ct))
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		// Суммарный размер всех загружаемых файлов не более 30 Мб.
		r.Body = http.MaxBytesReader(w, r.Body, maxMemory)
		r.ParseMultipartForm(maxMemory + 512)

		// Проверяем наличие файлов и строковых значений.
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

		// Получаем сформированную структуру заявки и возвращаем ошибку,
		// если code не равно 200.
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

		log.Debug("request body parsed succefully")

		// Получаем новый номер заявки и записываем его в структуру заявки.
		// ObjectID будет сгенерирован в методе БД AddReport. В случае ошибки
		// удаляем загруженные файлы из S3 хранилища.
		ctx := r.Context()
		newNum, err := st.CounterInc(ctx)
		if err != nil {
			go reports.RemoveFiles(l, report.Media, s3)
			log.Error("cannot receive new ID", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		report.Number = int64(newNum)

		// Добавляем сформированную заявку в БД. В случае ошибки удаляем
		// загруженные файлы из S3 хранилища.
		err = st.AddReport(ctx, report)
		if err != nil {
			go reports.RemoveFiles(l, report.Media, s3)
			log.Error("cannot add report to DB", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("new report added successfully")

		render.Status(r, http.StatusCreated)
		render.PlainText(w, r, strconv.Itoa(int(report.Number)))
	}
}
