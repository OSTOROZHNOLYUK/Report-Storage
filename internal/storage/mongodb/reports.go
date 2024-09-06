package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddReport добавляет одну новую заявку в БД.
func (s *Storage) AddReport(ctx context.Context, rep storage.Report) error {
	const operation = "storage.mongodb.AddReport"

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

// Reports возвращает все заявки из БД отсортированными по номеру
// в убывающем порядке.
func (s *Storage) Reports(ctx context.Context) ([]storage.Report, error) {
	const operation = "storage.mongodb.Reports"

	var reports []storage.Report
	collection := s.db.Database(dbName).Collection(colName)

	// Устанавливаем сортировку по полю number в убывающем порядке.
	opts := options.Find().SetSort(bson.D{{Key: "number", Value: -1}})

	// Получаем все заявки из БД.
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Записываем все заявки в массив структур.
	err = cursor.All(ctx, &reports)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if len(reports) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrNoData)
	}

	// Меняем местами долготу и широту.
	for i := range reports {
		reports[i].Geo.Coordinates[0], reports[i].Geo.Coordinates[1] = reports[i].Geo.Coordinates[1], reports[i].Geo.Coordinates[0]
	}

	return reports, nil
}
