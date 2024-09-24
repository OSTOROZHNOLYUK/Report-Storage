package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReportByNum возвращает заявку по ее уникальному номеру. Аргумент num
// должен быть больше 0, иначе вернет ошибку ErrIncorrectNum. Если документ
// с указанным номером не найден, то вернет ошибку ErrReportNotFound.
func (s *Storage) ReportByNum(ctx context.Context, num int) (storage.Report, error) {
	const operation = "storage.mongodb.ReportByNum"

	var report storage.Report
	if num < 1 {
		return report, fmt.Errorf("%s: %w", operation, storage.ErrIncorrectNum)
	}

	collection := s.db.Database(dbName).Collection(colReport)
	filter := bson.D{{Key: "number", Value: num}}
	err := collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return report, fmt.Errorf("%s: %w", operation, storage.ErrReportNotFound)
		}
		return report, fmt.Errorf("%s: %w", operation, err)
	}

	return report, nil
}
