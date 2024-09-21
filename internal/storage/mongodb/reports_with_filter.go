package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReportsWithFilter возвращает заявки в соответствии с переданными параметрами
// фильтра. Если параметр фильтра не задан или имеет некорректное значение, то
// используется значение по-умолчанию. Если заявки не найдены, то вернет ошибку
// ErrArrayNotFound.
func (s *Storage) ReportsWithFilter(ctx context.Context, fl storage.Filter) ([]storage.Report, error) {
	const operation = "storage.mongodb.ReportsWithFilter"

	var reports []storage.Report
	collection := s.db.Database(dbName).Collection(colName)

	// Задаем фильтр по статусам, если они переданы.
	filter := bson.D{}
	if len(fl.Status) > 0 {
		filter = bson.D{
			{Key: "status", Value: bson.M{"$in": fl.Status}},
		}
	}

	// Задаем порядок сортировки. По-умолчанию -1, нисходящий.
	sort := -1
	if fl.Sort == 1 {
		sort = 1
	}
	opts := options.Find().SetSort(bson.D{{Key: "number", Value: sort}})

	// Задаем количество. По-умолчанию 20.
	lim := 20
	if fl.Count > 0 {
		lim = fl.Count
	}
	opts.SetLimit(int64(lim))

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
