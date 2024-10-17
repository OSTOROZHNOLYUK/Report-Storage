package api

import (
	"Report-Storage/internal/storage"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// splitStatus преобразует строку с числами из query параметра
// в слайс статусов.
func splitStatus(s string) []storage.Status {
	var status []storage.Status
	if s == "" {
		return status
	}

	arr := strings.Split(s, ",")
	for _, v := range arr {
		n, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		status = append(status, storage.Status(n))
	}
	return status
}

// number получает значение параметра num из url запроса
// и преобразует в корректное значение.
func number(r *http.Request) (int, error) {
	numStr := chi.URLParam(r, "num")
	if numStr == "" {
		return 0, fmt.Errorf("empty number parameter")
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	if num < 1 {
		return 0, fmt.Errorf("number parameter less than 1")
	}
	return num, nil
}

// objectID получает значение параметра id из url запроса.
func objectID(r *http.Request) (string, error) {
	obj := chi.URLParam(r, "id")
	if obj == "" {
		return "", fmt.Errorf("empty objectID parameter")
	}
	return obj, nil
}

// point получает значения координат из query параметров x и y,
// и значение радиуса в метрах из query параметра r. Возвращает
// структуру точки и радиус.
func point(r *http.Request) (storage.Geo, int, error) {
	var position storage.Geo

	xParam := r.URL.Query().Get("x")
	yParam := r.URL.Query().Get("y")
	rParam := r.URL.Query().Get("r")

	if xParam == "" {
		return position, 0, fmt.Errorf("empty X parameter")
	}
	if yParam == "" {
		return position, 0, fmt.Errorf("empty Y parameter")
	}
	if rParam == "" {
		return position, 0, fmt.Errorf("empty R parameter")
	}

	x, err := strconv.ParseFloat(xParam, 64)
	if err != nil {
		return position, 0, fmt.Errorf("failed to parse X: %w", err)
	}
	y, err := strconv.ParseFloat(yParam, 64)
	if err != nil {
		return position, 0, fmt.Errorf("failed to parse Y: %w", err)
	}
	radius, err := strconv.Atoi(rParam)
	if err != nil {
		return position, 0, fmt.Errorf("failed to parse radius: %w", err)
	}

	position.Type = "Point"
	position.Coordinates = [2]float64{x, y}
	return position, radius, nil
}

// count получает значение из query параметра n и, если оно корректно,
// возвращает его. Иначе возвращает значение по умолчанию.
func count(r *http.Request) int {
	var def int = 20
	c := chi.URLParam(r, "n")
	if c == "" {
		return def
	}

	n, err := strconv.Atoi(c)
	if err != nil {
		return def
	}

	if n < 1 {
		return def
	}
	return n
}

// sort получает значение из query параметра sort и, если оно корректно,
// возвращает его. Иначе возвращает значение по умолчанию.
func sort(r *http.Request) int {
	var def int = -1
	s := chi.URLParam(r, "sort")
	if s == "" {
		return def
	}

	sort, err := strconv.Atoi(s)
	if err != nil {
		return def
	}

	if sort != 1 && sort != -1 {
		return def
	}
	return sort
}

// newStatus получает значение из query параметра new и, если оно корректно,
// возвращает его. Иначе возвращает ошибку.
func newStatus(r *http.Request) (storage.Status, error) {
	var status storage.Status
	s := r.URL.Query().Get("new")
	if s == "" {
		return 0, fmt.Errorf("empty new status parameter")
	}

	status = splitStatus(s)[0]
	if status < 1 || status > 5 {
		return 0, fmt.Errorf("incorrect new status value: %d", status)
	}
	return status, nil
}

// statusString возвращает строковое представление статуса заявки.
func statusString(s storage.Status) string {
	switch s {
	case 1:
		return "Неподтверждена"
	case 2:
		return "Создана"
	case 3:
		return "В работе"
	case 4:
		return "Завершена"
	case 5:
		return "Отклонена"
	default:
		return ""
	}
}
