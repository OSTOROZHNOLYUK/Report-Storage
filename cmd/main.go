package main

import (
	"Report-Storage/internal/config"
	"Report-Storage/internal/logger"
	"Report-Storage/internal/s3cloud"
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

	// Инициализируем клиент S3 хранилища.
	s3 := s3cloud.New(cfg.Endpoint, cfg.Bucket, cfg.AccessKey, cfg.SecretKey, cfg.Domain)
	log.Debug("S3 client initialized")

	// Инициализируем и запускаем HTTP сервер.
	srv := server.New(cfg)
	srv.API(log, st, s3)
	srv.Start()
	log.Info("Server started")

	// Блокируем выполнение основной горутины до сигнала прерывания.
	stopsignal.Stop()

	// После сигнала прерывания останавливаем сервер.
	srv.Shutdown()
	log.Info("Server stopped")
}
