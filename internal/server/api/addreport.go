package api

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/storage"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ReportAdder interface {
	AddReport(context.Context, storage.Report) error
}

func AddReport(log *slog.Logger, st ReportAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.api.AddReport"

		log.Info("new request to add report")

		// Ограничиваем чтение тела запроса размером в 1 мб.
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var rep storage.Report
		err := json.NewDecoder(r.Body).Decode(&rep)
		if err != nil {
			log.Error("cannot decode request body", logger.Err(err), slog.String("op", operation))
			http.Error(w, "incorrect report data", http.StatusBadRequest)
			return
		}
		slog.Debug("request body decoded")

		ctx := r.Context()
		err = st.AddReport(ctx, rep)
		if err != nil {
			log.Error("cannot add report to DB", logger.Err(err), slog.String("op", operation))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		slog.Debug("new report added successfully")

		w.WriteHeader(http.StatusCreated)
	}
}
