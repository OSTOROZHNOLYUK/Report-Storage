package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CounterInc увеличивает на единицу значение счетчика заявок в БД
// и возвращает это увеличенное значение.
func (s *Storage) CounterInc(ctx context.Context) (int32, error) {
	const operation = "storage.mongodb.Counter"

	collection := s.db.Database(dbName).Collection(colCounter)

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{{Key: "reportCount", Value: bson.D{{Key: "$exists", Value: true}}}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "reportCount", Value: 1},
		}},
	}

	var result bson.M
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		if c, ok := result["reportCount"].(int32); ok {
			return c, nil
		}
	}
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	if c, ok := result["reportCount"].(int32); ok {
		return c, nil
	}
	return 0, fmt.Errorf("%s: %w", operation, errors.New("incorrect type"))
}
