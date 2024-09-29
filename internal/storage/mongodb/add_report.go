package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddReport добавляет одну новую заявку в БД. Не осуществляет валидацию
// rep, ожидает полностью валидную заявку.
func (s *Storage) AddReport(ctx context.Context, rep storage.Report) error {
	const operation = "storage.mongodb.AddReport"

	if rep.Number < 1 {
		return fmt.Errorf("%s: %w", operation, storage.ErrIncorrectNum)
	}

	// Устанавливаем ObjectID.
	rep.ID = primitive.NewObjectID()

	// Меняем местами широту и долготу.
	rep.Geo.Coordinates[0], rep.Geo.Coordinates[1] = rep.Geo.Coordinates[1], rep.Geo.Coordinates[0]

	// Производим вставку новой заявки.
	collection := s.db.Database(dbName).Collection(colReport)
	_, err := collection.InsertOne(ctx, rep)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
