// Пакет для работы с базой данных MongoDB.

package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	testDatabase   = "unitTestDB"
	testCollection = "unitTestCollection"
	testCounter    = "unitTestCounter"
)

// path - адрес БД для юнит-тестов.
var path string = "mongodb://194.54.157.224:10501/"

// reports - заявки для юнит-тестов.
var reports = []storage.Report{
	{
		Number:      1,
		Address:     "Адрес 1",
		Description: "Описание заявки 1",
		Contacts:    storage.Contacts{Email: "bob@gmail.com", Telegram: "@bob"},
		Media:       []string{"https://google.com"},
		Geo:         storage.Geo{Coordinates: [2]float64{55.75388130172051, 37.62026781374883}},
	},
	{
		Number:      2,
		Address:     "Адрес 2",
		Description: "Описание заявки 2",
		Contacts:    storage.Contacts{Email: "bill@gmail.com", Whatsapp: "+71234567890"},
		Media:       []string{"https://google.com"},
		Geo:         storage.Geo{Coordinates: [2]float64{55.75909434896026, 37.619124583054855}},
	},
	{
		Number:      3,
		Address:     "Адрес 3",
		Description: "Описание заявки 3",
		Media:       []string{"https://google.com"},
		Geo:         storage.Geo{Coordinates: [2]float64{59.939543808173305, 30.31511987692599}},
	},
}

// addOne добавляет одну заявку в БД. Функция для использования в тестах.
func (s *Storage) addOne(rep storage.Report) (string, error) {

	rep.ID = primitive.NewObjectID()
	rep.Created = time.Now()
	rep.Updated = time.Now()
	rep.Status = storage.Unverified
	rep.Geo.Type = "Point"
	rep.Geo.Coordinates[0], rep.Geo.Coordinates[1] = rep.Geo.Coordinates[1], rep.Geo.Coordinates[0]

	collection := s.db.Database(testDatabase).Collection(testCollection)
	res, err := collection.InsertOne(context.Background(), rep)
	if err != nil {
		return "", err
	}
	hex := res.InsertedID.(primitive.ObjectID)
	return hex.Hex(), nil
}

// getOne возвращает заявку по ObjectID. Функция для использования в тестах.
func (s *Storage) getOne(id string) (storage.Report, error) {

	var report storage.Report
	collection := s.db.Database(testDatabase).Collection(testCollection)

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return report, err
	}
	filter := bson.D{{Key: "_id", Value: obj}}
	err = collection.FindOne(context.Background(), filter).Decode(&report)
	if err != nil {
		return report, err
	}
	return report, nil
}

// trun удаляет все записи в колекции. Функция для использования в тестах.
func (s *Storage) trun(name string) error {
	collection := s.db.Database(dbName).Collection(name)
	// err := collection.Drop(context.Background())
	_, err := collection.DeleteMany(context.Background(), bson.D{})
	return err
}

func Test_new(t *testing.T) {
	opts := setOpts(path, "admin", os.Getenv("MONGO_DB_PASSWD"))

	st, err := new(opts)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
}

func Test_checkStatus(t *testing.T) {

	tests := []struct {
		name   string
		status storage.Status
		want   bool
	}{
		{
			name:   "OK",
			status: 1,
			want:   true,
		},
		{
			name:   "Error",
			status: 6,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkStatus(tt.status); got != tt.want {
				t.Errorf("checkStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
