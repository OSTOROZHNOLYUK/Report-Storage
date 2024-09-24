package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReportByID возвращает заявку по ее ObjectID. Аргумент id должен быть
// валидной hex строкой. Если id пустая строка или имеет некорректный
// формат, то вернет ошибку ErrIncorrectID. Если документ с указанным
// id не найден, то вернет ошибку ErrReportNotFound.
func (s *Storage) ReportByID(ctx context.Context, id string) (storage.Report, error) {
	const operation = "storage.mongodb.ReportByID"

	var report storage.Report
	if id == "" {
		return report, fmt.Errorf("%s: %w", operation, storage.ErrIncorrectID)
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return report, fmt.Errorf("%s: %w", operation, storage.ErrIncorrectID)
	}

	collection := s.db.Database(dbName).Collection(colReport)
	filter := bson.D{{Key: "_id", Value: obj}}
	err = collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return report, fmt.Errorf("%s: %w", operation, storage.ErrReportNotFound)
		}
		return report, fmt.Errorf("%s: %w", operation, err)
	}

	return report, nil
}
