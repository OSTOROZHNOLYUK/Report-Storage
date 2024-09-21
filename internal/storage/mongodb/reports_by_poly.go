package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// ReportsByPoly возвращает все заявки в границах многоугольника с фильтрацией
// по статусам. Если в параметр status передать nil или пустой слайс, то вернет
// все заявки. Параметр poly - слайс массивов по два элемента, представляет
// список точек координат по периметру многоугольника. ReportsByPoly не проверяет
// принимаемые аргументы, ожидает полностью валидные значения. Если заявки
// не найдены, то вернет ошибку ErrArrayNotFound.
func (s *Storage) ReportsByPoly(ctx context.Context, poly [][2]float64, status []storage.Status) ([]storage.Report, error) {
	const operation = "storage.mongodb.ReportsByPoly"

	var reports []storage.Report
	collection := s.db.Database(dbName).Collection(colName)

	// Меняем местами широту и долготу.
	for k := range poly {
		poly[k][0], poly[k][1] = poly[k][1], poly[k][0]
	}
	// Добавляем первую точку координат в конец слайса, замыкая многоугольник.
	last := poly[0]
	poly = append(poly, last)
	// Добавляем одну размерность и формируем GeoJSON.
	var pp = [][][2]float64{poly}
	polygon := bson.D{{Key: "type", Value: "Polygon"}, {Key: "coordinates", Value: pp}}

	// Создаем фильтр из многоугольника.
	filter := bson.M{}
	filter["geo"] = bson.D{
		{Key: "$geoWithin", Value: bson.D{
			{Key: "$geometry", Value: polygon},
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
