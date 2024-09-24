package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// DeleteRejected удаляет все заявки со статусом Rejected. Возвращает
// количество удаленных документов и ошибку. Если не было удалено
// ни одного документа, то вернет ошибку ErrReportNotFound.
func (s *Storage) DeleteRejected(ctx context.Context) (int, error) {
	const operation = "storage.mongodb.DeleteRejected"

	collection := s.db.Database(dbName).Collection(colReport)
	filter := bson.D{{Key: "status", Value: storage.Rejected}}
	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	if res.DeletedCount == 0 {
		return 0, fmt.Errorf("%s: %w", operation, storage.ErrReportNotFound)
	}
	return int(res.DeletedCount), nil
}
