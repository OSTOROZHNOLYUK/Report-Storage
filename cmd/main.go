package main

import (
	"Report-Storage/internal/config"
	"Report-Storage/internal/logger"
	"Report-Storage/internal/server"
	"Report-Storage/internal/stopsignal"
	"Report-Storage/internal/storage/mongodb"
)

func main() {

	// Инициализируем конфиг и логгер.
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Debug("Config file and logger initialized")

	// Инициализируем пул подключений БД.
	st := mongodb.New(cfg)
	log.Debug("Storage initialized")

	// Инициализируем и запускаем HTTP сервер.
	srv := server.New(cfg)
	srv.Middleware()
	srv.API(log, st)
	srv.Start()
	log.Info("Server started")

	// Блокируем выполнение основной горутины до сигнала прерывания.
	stopsignal.Stop()

	// После сигнала прерывания останавливаем сервер.
	srv.Shutdown()
	log.Info("Server stopped")
}
