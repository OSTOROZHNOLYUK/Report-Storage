package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateStatus изменяет значение статуса у заявки по переданному номеру.
// Аргумент num должен быть больше 0, иначе вернет ошибку ErrIncorrectNum.
// Аргумент status должен быть валидным значением статуса, иначе вернет
// ошибку ErrIncorrectStatus. Если документ с указанным номером не найден,
// то вернет ошибку ErrReportNotFound.
func (s *Storage) UpdateStatus(ctx context.Context, num int, status storage.Status) (storage.Report, error) {
	const operation = "storage.mongodb.UpdateStatus"

	var report storage.Report
	if num < 1 {
		return report, fmt.Errorf("%s: %w", operation, storage.ErrIncorrectNum)
	}
	if !checkStatus(status) {
		return report, fmt.Errorf("%s: %w", operation, storage.ErrIncorrectStatus)
	}

	collection := s.db.Database(dbName).Collection(colReport)
	filter := bson.D{{Key: "number", Value: num}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: status},
			{Key: "updated", Value: time.Now()},
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&report)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return report, fmt.Errorf("%s: %w", operation, storage.ErrReportNotFound)
		}
		return report, fmt.Errorf("%s: %w", operation, err)
	}

	return report, nil
}
