package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Reports возвращает все заявки из БД отсортированными по номеру
// в убывающем порядке.
func (s *Storage) Reports(ctx context.Context, status []storage.Status) ([]storage.Report, error) {
	const operation = "storage.mongodb.Reports"

	var reports []storage.Report
	collection := s.db.Database(dbName).Collection(colName)

	filter := bson.D{}

	// TODO: фильтр по статусам.
	//
	// if status != nil && len(status) > 0 {
	// }

	// Устанавливаем сортировку по полю number в убывающем порядке.
	opts := options.Find().SetSort(bson.D{{Key: "number", Value: -1}})

	// Получаем все заявки из БД.
	cursor, err := collection.Find(ctx, filter, opts)
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
