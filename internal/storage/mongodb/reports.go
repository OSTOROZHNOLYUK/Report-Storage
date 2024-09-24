package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Reports возвращает заявки из БД отсортированными по номеру в убывающем
// порядке. Вторым параметром принимает слайс статусов и возвращает все
// заявки с указанными статусами. Если передать nil или пустой слайс, то
// вернет все заявки. Если заявки не найдены, то вернет ошибку ErrArrayNotFound.
func (s *Storage) Reports(ctx context.Context, status []storage.Status) ([]storage.Report, error) {
	const operation = "storage.mongodb.Reports"

	var reports []storage.Report
	collection := s.db.Database(dbName).Collection(colReport)

	// Задаем фильтр по статусам, если они переданы.
	filter := bson.D{}
	if len(status) > 0 {
		filter = bson.D{
			{Key: "status", Value: bson.M{"$in": status}},
		}
	}

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
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrArrayNotFound)
	}

	// Меняем местами долготу и широту.
	for i := range reports {
		reports[i].Geo.Coordinates[0], reports[i].Geo.Coordinates[1] = reports[i].Geo.Coordinates[1], reports[i].Geo.Coordinates[0]
	}

	return reports, nil
}
