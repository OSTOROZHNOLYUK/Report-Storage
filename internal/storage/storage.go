// Пакет storage содержит структуры для работы с базой данных, возможные
// возвращаемые ошибки и интерфейс подключения к БД.
package storage

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrIncorrectNum    = errors.New("incorrect report number")
	ErrIncorrectID     = errors.New("incorrect report objectid")
	ErrIncorrectStatus = errors.New("incorrect report status")
	ErrReportNotFound  = errors.New("report not found")
	ErrArrayNotFound   = errors.New("reports array not found")
)

// Status - целочисленное выражение статуса заявки.
type Status int

// Константы статусов заявки.
const (
	Unverified Status = 1 + iota
	Opened
	InProgress
	Closed
	Rejected
	StatusPending
	StatusActive
	StatusCompleted
)

// Geo - тип данных географических координат точки.
type Geo struct {
	// Type - тип объекта, в нашем случае всегда значение "Point".
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
	Status Status `json:"status" bson:"status"`
}

// Filter - структура фильтра для получения заявок.
type Filter struct {
	// Count отражает необходимое количество заявок, должно быть > 0.
	Count int
	// Sort указывает порядок сортировки по номеру, значение должно
	// быть 1 для восходящего и -1 для нисходящего порядков.
	Sort int
	// Слайс статусов.
	Status []Status
}

// Statistic - структура статистики заявок со статусами.
type Statistic struct {
	Total, Unverified, Opened, InProgress, Closed, Rejected int
}

// StatusFromString преобразует строку в тип Status.
func StatusFromString(s string) (Status, error) {
	switch s {
	case "pending":
		return StatusPending, nil
	case "active":
		return StatusActive, nil
	case "completed":
		return StatusCompleted, nil
	// обработайте другие случаи
	default:
		return -1, fmt.Errorf("неизвестный статус: %s", s)
	}
}
