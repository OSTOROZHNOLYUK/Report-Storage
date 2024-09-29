package s3cloud

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// FileStorage - структура клиента S3 хранилища.
type FileStorage struct {
	storage  *minio.Client
	endpoint string
	bucket   string
	domain   string
}

// UploadInput - структура файла для загрузки в хранилище.
type UploadInput struct {
	File        io.Reader
	Name        string
	Size        int64
	ContentType string
}

// New - конструктор клиента S3 хранилища.
func New(endpoint, bucket, accessKey, secretKey, domain string) *FileStorage {
	s3, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("failed to create S3 storage client: %s", err.Error())
	}

	fs := &FileStorage{
		storage:  s3,
		endpoint: endpoint,
		bucket:   bucket,
		domain:   domain,
	}
	return fs
}

// Upload загружает файл в хранилище. Возвращает url загруженного файла,
// либо ошибку.
func (fs *FileStorage) Upload(ctx context.Context, input UploadInput) (string, error) {
	const operation = "s3cloud.Upload"

	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	_, err := fs.storage.PutObject(ctx, fs.bucket, input.Name, input.File, input.Size, opts)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	url := fmt.Sprintf("%s/%s", fs.domain, input.Name)
	return url, nil
}

// Remove удаляет файл с переданным url из хранилища.
func (fs *FileStorage) Remove(ctx context.Context, url string) error {
	const operation = "s3cloud.Remove"

	str := strings.Split(url, "/")
	name := str[len(str)-1]

	opts := minio.RemoveObjectOptions{GovernanceBypass: true}
	err := fs.storage.RemoveObject(ctx, fs.bucket, name, opts)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
