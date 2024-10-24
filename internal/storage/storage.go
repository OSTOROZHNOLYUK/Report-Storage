// Пакет storage содержит структуры для работы с базой данных, возможные
// возвращаемые ошибки и интерфейс подключения к БД.
package storage

import (
	"errors"
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
)

// Geo - тип данных географических координат точки.
type Geo struct {
	// Type - тип объекта, в нашем случае всегда значение "Point".
	Type string `json:"type,omitempty" bson:"type" validate:"omitempty,max=100"`
	// Coordinates - координаты, первый элемент - широта, второй элемент - долгота.
	Coordinates [2]float64 `json:"coordinates" bson:"coordinates" validate:"required,dive,required"`
}

// Contacts - структура контактов отправителя заявки.
type Contacts struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty" validate:"omitempty,email,max=100"`
	Whatsapp string `json:"whatsapp,omitempty" bson:"whatsapp,omitempty" validate:"omitempty,max=100"` // e164
	Telegram string `json:"telegram,omitempty" bson:"telegram,omitempty" validate:"omitempty,max=100"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty" validate:"omitempty,max=100"` // e164
}

// Report - основная структура заявки о проблеме.
type Report struct {
	// ID хранит значение ObjectID, используемое в MongoDB.
	ID primitive.ObjectID `json:"id" bson:"_id" validate:"required,max=100"`

	// Number содержит уникальный порядковый номер заявки.
	Number int64 `json:"number" bson:"number" validate:"required,number,min=1,max=99999999"`

	// Created содержит время создания заявки с БД.
	Created time.Time `json:"created" bson:"created" validate:"required"`

	// Updated содержит время последнего изменения заявки.
	Updated time.Time `json:"updated,omitempty" bson:"updated"`

	// City содержит значение города или местности заявки.
	City string `json:"city" bson:"city" validate:"required,max=100"`

	// Address хранит строковое представление ближайшего адреса.
	Address string `json:"address" bson:"address" validate:"required,max=100"`

	// Description содержит описание заявки клиентом в свободной форме.
	Description string `json:"description" bson:"description" validate:"omitempty,max=300"`

	// Contacts содержит возможные контакты клиента.
	Contacts Contacts `json:"contacts,omitempty" bson:"contacts,omitempty"`

	// Media содержит слайс ссылок на медиа файлы по заявке.
	Media []string `json:"media" bson:"media" validate:"required,min=1,max=5"`

	// Тип Coordinates хранит географические координаты заявки.
	Geo Geo `json:"geo" bson:"geo" validate:"required"`

	// Status содержит целочисленную константу, отражающую текущий
	// статус заявки.
	Status Status `json:"status" bson:"status" validate:"required,number,min=1,max=5"`
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
// func StatusFromString(s string) (Status, error) {
// 	switch s {
// 	case "pending":
// 		return StatusPending, nil
// 	case "active":
// 		return StatusActive, nil
// 	case "completed":
// 		return StatusCompleted, nil
// 	// обработайте другие случаи
// 	default:
// 		return -1, fmt.Errorf("неизвестный статус: %s", s)
// 	}
// }
