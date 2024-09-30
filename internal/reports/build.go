package reports

import (
	"Report-Storage/internal/logger"
	"Report-Storage/internal/s3cloud"
	"Report-Storage/internal/storage"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/h2non/filetype/matchers"
)

const (
	// jsonInputName - имя поля ввода на строне клиента.
	jsonInputName string = "json"
	// maxFile - максимальный размер одного файла.
	maxFile int64 = 5 << 20
	// maxPic - максимальная длина стороны фото.
	maxPic int = 1800
	// jpegQuality - уровень качества формата JPEG при кодировании.
	jpegQuality int = 60
)

// Request - структура запроса на добавление новой заявки.
type Request struct {
	City        string           `json:"city" validate:"required,max=100"`
	Address     string           `json:"address" validate:"required,max=100"`
	Description string           `json:"description,omitempty" validate:"max=300"`
	Contacts    storage.Contacts `json:"contacts,omitempty" validate:"omitempty"`
	Geo         storage.Geo      `json:"geo" validate:"required"`
}

// ReportAdder - интерфейс для БД в обработчике AddReport.
type ReportAdder interface {
	AddReport(context.Context, storage.Report) error
	CounterInc(context.Context) (int32, error)
}

// FileSaver - интерфейс для объектного хранилища в обработчике AddReport.
type FileSaver interface {
	Upload(context.Context, s3cloud.UploadInput) (string, error)
	Remove(context.Context, string) error
}

// Build формирует структуру заявки storage.Report из multipart запроса.
// Этот запрос должен содержать часть с именем "json", где передается
// JSON новой заявки, и от 1 до 5 частей с любыми именами, содержащими
// файл в формате jpeg или png. Любые другие строковые части игнорируются,
// любые другие файлы вернут ошибку на запрос.
// Часть json распарсивается в структуру Request, файлы перекодируются в
// jpeg с заданным качеством и загружаются в объектное хранилище.
// Функция возвращает структуру заявки и HTTP код как символ ошибки. Если
// код не равен 200, то при обработке возникли ошибки, и структура заявки
// будет пуста.
func Build(l *slog.Logger, s3 FileSaver, r *http.Request) (storage.Report, int) {
	const operation = "reports.Build"

	log := l.With(
		slog.String("op", operation),
	)

	var req Request
	var report storage.Report
	var code int
	ctx := r.Context()

	// Обработка JSON.

	// Вычитываем JSON из запроса в структуру Request.
	for key, body := range r.MultipartForm.Value {
		if key != jsonInputName {
			continue
		}
		for _, value := range body {
			err := render.DecodeJSON(strings.NewReader(value), &req)
			if err != nil {
				log.Error("cannot decode request json", logger.Err(err))
				return report, http.StatusBadRequest
			}
		}
	}

	// Валидируем поля запроса.
	valid := validator.New()
	err := valid.Struct(req)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("validation failed", logger.Err(validateErr))
		return report, http.StatusBadRequest
	}
	log.Debug("json input decoded and validated successfully")

	// Обработка файлов.

	// Обрабатываем каждый файл в отдельной горутине. Результаты пишем
	// в канал urls, ошибки в канал errFiles.
	var wg sync.WaitGroup
	urls := make(chan string, len(r.MultipartForm.File))
	errFiles := make(chan error, len(r.MultipartForm.File))

	for _, body := range r.MultipartForm.File {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, part := range body {

				// Проверяем размер файла, не больше maxFile (5 Мб).
				if part.Size > maxFile {
					log.Error(
						"file size is greater than maxFile value",
						slog.String("filename", part.Filename),
						slog.Int64("size", part.Size),
					)
					code = http.StatusBadRequest
					errFiles <- errors.New("file size too large")
					return
				}

				// Открываем файл и считываем его в буфер.
				file, err := part.Open()
				if err != nil {
					log.Error(
						"cannot open multipart file data",
						slog.String("filename", part.Filename),
						logger.Err(err),
					)
					code = http.StatusBadRequest
					errFiles <- err
					return
				}
				defer file.Close()

				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(file)
				if err != nil {
					log.Error(
						"cannot read file data to buffer",
						slog.String("filename", part.Filename),
						logger.Err(err),
					)
					code = http.StatusBadRequest
					errFiles <- err
					return
				}

				// Проверяем, что файл имеет тип jpeg или png.
				b := buf.Bytes()
				if !matchers.Jpeg(b) && !matchers.Png(b) {
					log.Error("unsupported file type", slog.String("filename", part.Filename))
					code = http.StatusUnsupportedMediaType
					errFiles <- errors.New("unsupported file type")
					return
				}

				// Декодируем файл в image.Image и определяем ориентацию фото.
				var h, w int
				img, err := imaging.Decode(buf)
				if err != nil {
					log.Error(
						"failed to decode file to image.Image",
						slog.String("filename", part.Filename),
						logger.Err(err),
					)
					code = http.StatusInternalServerError
					errFiles <- err
				}
				if img.Bounds().Dx() > img.Bounds().Dy() {
					w = maxPic
				} else {
					h = maxPic
				}

				// Меняем размер фото на максимально допустимый, затем кодируем
				// фото в JPEG формат с заданным качеством.
				dstImage := imaging.Resize(img, w, h, imaging.Lanczos)
				buf = new(bytes.Buffer)
				opts := imaging.JPEGQuality(jpegQuality)
				err = imaging.Encode(buf, dstImage, imaging.JPEG, opts)
				if err != nil {
					log.Error(
						"failed to encode image to jpeg",
						slog.String("filename", part.Filename),
						logger.Err(err),
					)
					code = http.StatusInternalServerError
					errFiles <- err
				}
				fReader := bytes.NewReader(buf.Bytes())

				// Создаем структуру для загрузки файла в S3 хранилище и загружаем ее.
				// Полученную ссылку на файл отправляем в канал urls.
				input := s3cloud.UploadInput{
					File:        fReader,
					Name:        generateFileNameJPEG(),
					Size:        fReader.Size(),
					ContentType: part.Header.Get("Content-Type"),
				}
				url, err := s3.Upload(ctx, input)
				if err != nil {
					log.Error(
						"cannot upload file to s3",
						slog.String("filename", part.Filename),
						logger.Err(err),
					)
					code = http.StatusInternalServerError
					errFiles <- err
					return
				}
				urls <- url
			}
		}()
	}

	// Ждем завершения обработки и загрузки всех файлов.
	wg.Wait()

	close(urls)
	close(errFiles)

	// Вычитываем ссылки из канала в слайс.
	var media []string
	for v := range urls {
		media = append(media, v)
	}

	// Если при обработке файлов возникли ошибки, то асинхронно удаляем
	// успешно загруженные файлы из S3, так как весь запрос должен завершиться
	// ошибкой.
	if len(errFiles) > 0 {
		go RemoveFiles(log, media, s3)
		log.Error("failed to upload some files")
		return report, code
	}

	// Формируем все поля структуры заявки кроме ID и Number. Эти поля будут
	// заполнены значениями на других уровнях.
	report.Created = time.Now()
	report.Updated = time.Now()
	report.City = req.City
	report.Address = req.Address
	report.Description = req.Description
	report.Contacts = req.Contacts
	report.Media = media
	report.Geo = req.Geo
	report.Geo.Type = "Point"
	report.Status = storage.Unverified

	// Возвращаем валидную заявку и код 200.
	return report, http.StatusOK
}
