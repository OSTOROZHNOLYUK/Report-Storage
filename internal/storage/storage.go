// Пакет storage содержит структуры для работы с базой данных, возможные
// возвращаемые ошибки и интерфейс подключения к БД.
package storage

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoData = errors.New("no data")
)

// Целочисленные константы статусов заявки.
const (
	Unverified = 1 + iota
	Opened
	InProgress
	Closed
	Rejected
)

// Geo - тип данных географических координат точки.
type Geo struct {
	// Type - тип объекта, в нашем случае всегда Point.
	Type string `json:"type" bson:"type"`
	// Coordinates - координаты, первый элемент - широта, второй элемент - долгота.
	Coordinates [2]float64 `json:"coordinates" bson:"coordinates"`
}

// Contacts - структура контактов отправителя заявки.
type Contacts struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Whatsapp string `json:"whatsapp,omitempty" bson:"whatsapp,omitempty"`
	Telegram string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
}

// Report - основная структура заявки о проблеме.
type Report struct {
	// ID хранит значение ObjectID, используемое в MongoDB.
	ID primitive.ObjectID `json:"id" bson:"_id"`

	// Number содержит номер заявки, сгенерированный при ее создании в сервисе
	// обработки новых заявок.
	Number int64 `json:"number" bson:"number"`

	// Created содержит время создания заявки с БД.
	Created time.Time `json:"created" bson:"created"`

	// Updated содержит время последнего изменения заявки.
	Updated time.Time `json:"updated" bson:"updated"`

	// Address хранит строковое представление ближайшего адреса. Указывается
	// клиентом при создании заявки.
	Address string `json:"address" bson:"address"`

	// Description содержит описание заявки клиентом в свободной форме.
	Description string `json:"description" bson:"description"`

	// Contacts содержит возможные контакты клиента.
	Contacts Contacts `json:"contacts,omitempty" bson:"contacts,omitempty"`

	// Media содержит срез ссылок на медиа файлы по заявке.
	Media []string `json:"media" bson:"media"`

	// Тип Coordinates хранит географические координаты заявки.
	Geo Geo `json:"geo" bson:"geo"`

	// Status содержит целочисленную константу, отражающую текущий
	// статус заявки.
	Status int `json:"status" bson:"status"`
}
