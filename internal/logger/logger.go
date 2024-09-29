// Пакет logger используется для работы с логгером из пакета slog
// стандартной библиотеки.
package logger

import (
	"io"
	"log/slog"
	"os"
)

// Константы окружений
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// SetupLOgger инициализирует логгер из пакета slog с выводом в текстовом
// формате либо JSON в зависимости от окружения из конфиг файла.
// Default значение - prod.
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

// SetupDiscard инициализирует логгер из пакета slog, который никуда не пишет.
// Функция для использования в тестах.
func SetupDiscard() *slog.Logger {
	log := slog.New(
		slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}),
	)
	return log
}

// Err - обертка для ошибки, представляет ее как атрибут слоггера.
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
