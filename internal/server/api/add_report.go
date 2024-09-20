package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

// ReportAdder - интерфейс для БД в обработчике AddReport.
type ReportAdder interface {
	AddReport(context.Context, storage.Report) error
}

// AddReport обрабатывает запрос на добавление новой заявки в хранилище.
// При успехе возвращает код 201.
func AddReport(l *slog.Logger, st ReportAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.AddReport"

		log := l.With(
			slog.String("op", operation),
		)
		log.Info("request to add new report")

		// Ограничиваем чтение тела запроса размером в 1 мб.
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var rep storage.Report
		err := render.DecodeJSON(r.Body, &rep)
		if err != nil {
			log.Error("cannot decode request body", logger.Err(err))
			http.Error(w, "incorrect report data", http.StatusBadRequest)
			return
		}
		log.Debug("request body decoded")

		// TODO: валидация полей заявки.

		ctx := r.Context()
		err = st.AddReport(ctx, rep)
		if err != nil {
			log.Error("cannot add report to DB", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Debug("new report added successfully")

		render.Status(r, http.StatusCreated)
	}
}
