package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddReport добавляет одну новую заявку в БД. Не осуществляет валидацию
// rep, ожидает полностью валидную заявку.
func (s *Storage) AddReport(ctx context.Context, rep storage.Report) error {
	const operation = "storage.mongodb.AddReport"

	// TODO: вынести всю валидацию полей rep на уровень выше.

	if rep.Number < 1 {
		return fmt.Errorf("%s: %w", operation, storage.ErrIncorrectNum)
	}

	// Устанавливаем ObjectID, время создания, изменения и статус.
	rep.ID = primitive.NewObjectID()
	rep.Created = time.Now()
	rep.Updated = time.Now()
	rep.Status = storage.Unverified

	// Добавляем тип объекта и меняем местами широту и долготу.
	rep.Geo.Type = "Point"
	rep.Geo.Coordinates[0], rep.Geo.Coordinates[1] = rep.Geo.Coordinates[1], rep.Geo.Coordinates[0]

	// Производим вставку новой заявки.
	collection := s.db.Database(dbName).Collection(colName)
	_, err := collection.InsertOne(ctx, rep)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
