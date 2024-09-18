package api

import (
	"Report-Storage/internal/storage"
	"strconv"
	"strings"
)

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
