package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// DeleteByNum удаляет заявку по ее уникальному номеру. Аргумент num
// должен быть больше 0, иначе вернет ошибку ErrIncorrectNum. Если документ
// с указанным номером не найден, то вернет ошибку ErrReportNotFound.
func (s *Storage) DeleteByNum(ctx context.Context, num int) error {
	const operation = "storage.mongodb.DeleteByNum"

	if num < 1 {
		return fmt.Errorf("%s: %w", operation, storage.ErrIncorrectNum)
	}

	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "number", Value: num}}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", operation, storage.ErrReportNotFound)
	}
	return nil
}
