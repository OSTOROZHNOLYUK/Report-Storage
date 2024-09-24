package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Statistic возвращает общее количество заявок и отдельно по статусам.
// Если в коллекции нет документов, то вернет 0 по всем статусам и nil.
func (s *Storage) Statistic(ctx context.Context) (storage.Statistic, error) {
	const operation = "storage.mongodb.Statistic"

	var stat storage.Statistic
	collection := s.db.Database(dbName).Collection(colName)

	// Получаем общее количество заявок.
	c, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return stat, fmt.Errorf("%s: %w", operation, err)
	}
	if c == 0 {
		return stat, nil
	}
	stat.Total = int(c)

	// Создаем агрегацию для подсчета количества заявок по статусам.
	group := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{
				{Key: "$count", Value: bson.D{}},
			}},
		}},
	}
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{group})
	if err != nil {
		return stat, fmt.Errorf("%s: %w", operation, err)
	}

	var results []bson.M
	err = cursor.All(ctx, &results)
	if err != nil {
		return stat, fmt.Errorf("%s: %w", operation, err)
	}
	fmt.Println(results)
	for _, r := range results {
		status := r["_id"].(int32)
		switch status {
		case 1:
			stat.Unverified = int(r["count"].(int32))
		case 2:
			stat.Opened = int(r["count"].(int32))
		case 3:
			stat.InProgress = int(r["count"].(int32))
		case 4:
			stat.Closed = int(r["count"].(int32))
		case 5:
			stat.Rejected = int(r["count"].(int32))
		}
	}

	return stat, nil
}
