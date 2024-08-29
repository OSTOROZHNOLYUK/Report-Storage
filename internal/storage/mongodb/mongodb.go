// Пакет для работы с базой данных MongoDB.
package mongodb

import (
	"Report-Storage/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Название базы и коллекции в БД. Используются переменные вместо констант,
// так как в тестах им присваиваются другие значения.
var (
	dbName  string = "reportStorage"
	colName string = "reports"
)

// tmConn - таймаут на создание пула подключений.
const tmConn time.Duration = time.Second * 10

// Storage - пул подключений к БД.
type Storage struct {
	db *mongo.Client
}

// New - обертка для конструктора пула подключений new.
func New(cfg *config.Config) *Storage {
	opts := setOptsNoPasswd(cfg.StoragePath)
	storage, err := new(opts)
	if err != nil {
		log.Fatalf("failed to init storage: %s", err.Error())
	}
	return storage
}

// setOpts возвращает опции нового подключения с авторизацией.
// Пока закомментировал, чтобы IDE не ругалась.
//
// func setOpts(path, user, password string) *options.ClientOptions {
// 	credential := options.Credential{
// 		AuthMechanism: "SCRAM-SHA-256",
// 		AuthSource:    "admin",
// 		Username:      user,
// 		Password:      password,
// 	}
// 	opts := options.Client().ApplyURI(path).SetAuth(credential)
// 	return opts
// }

// setOptsNoPasswd возвращает опции нового подключения без авторизации.
func setOptsNoPasswd(path string) *options.ClientOptions {
	return options.Client().ApplyURI(path)
}

// new - конструктор пула подключений к БД.
func new(opts *options.ClientOptions) (*Storage, error) {
	const operation = "storage.mongodb.new"

	tm, cancel := context.WithTimeout(context.Background(), tmConn)
	defer cancel()

	db, err := mongo.Connect(tm, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

// Close - обертка для закрытия пула подключений.
func (s *Storage) Close() error {
	return s.db.Disconnect(context.Background())
}
