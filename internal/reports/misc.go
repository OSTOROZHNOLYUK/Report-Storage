package reports

import (
	"Report-Storage/internal/logger"
	"context"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/sqids/sqids-go"
)

// generateFileName генерирует имя для файла в виде строки из
// закодированного текущего времени, случайного числа от 1 до 9999
// и расширения файла.
func generateFileNameJPEG() string {
	tm := time.Now()
	sec := uint64(tm.Unix())
	nano := uint64(tm.Nanosecond())

	s, err := sqids.New(sqids.Options{MinLength: 12})
	if err != nil {
		return ""
	}
	name, err := s.Encode([]uint64{sec, nano})
	if err != nil {
		return ""
	}
	suff := strconv.Itoa(rand.Intn(9999))
	ext := ".jpg"

	return strings.Join([]string{name, suff, ext}, "")
}

// RemoveFiles удаляет все файлы из S3 хранилища по url из переданного слайса.
func RemoveFiles(log *slog.Logger, urls []string, s3 FileSaver) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	for _, url := range urls {
		err := s3.Remove(ctx, url)
		if err != nil {
			log.Error("failed to remove file from S3", slog.String("url", url), logger.Err(err))
		}
	}
}
