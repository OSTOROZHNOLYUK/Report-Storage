package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// ReportsByRadius возвращает все заявки в радиусе от точки с фильтрацией
// по статусам. Если в параметр status передать nil или пустой слайс, то
// вернет все заявки. r - радиус в метрах; p - структура точки координат
// storage.Geo, где поле Type должно иметь значение "Point". Не проверяет
// принимаемые аргументы, ожидает полностью валидные значения. Если заявки
// не найдены, то вернет ошибку ErrArrayNotFound.
func (s *Storage) ReportsByRadius(ctx context.Context, r int, p storage.Geo, status []storage.Status) ([]storage.Report, error) {
	const operation = "storage.mongodb.ReportsByRadius"

	var reports []storage.Report
	collection := s.db.Database(dbName).Collection(colReport)

	// Меняем местами широту и долготу, затем формируем GeoJSON.
	p.Coordinates[0], p.Coordinates[1] = p.Coordinates[1], p.Coordinates[0]
	point := bson.D{{Key: "type", Value: p.Type}, {Key: "coordinates", Value: p.Coordinates}}

	// Создаем фильтр из точки и радиуса.
	filter := bson.M{}
	filter["geo"] = bson.D{
		{Key: "$near", Value: bson.D{
			{Key: "$geometry", Value: point},
			{Key: "$maxDistance", Value: r},
		}},
	}

	// Расширяем фильтр статусами, если они переданы.
	if len(status) > 0 {
		filter["status"] = bson.M{"$in": status}
	}

	// Получаем все заявки из БД.
	cursor, err := collection.Find(ctx, filter)
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
