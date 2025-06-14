// Package s3 предоставляет функционал для работы с AWS S3 хранилищем
package s3

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ErrorFilesCountMismatch возникает когда количество файлов не соответствует количеству ключей
var ErrorFilesCountMismatch = errors.New("количество файлов не совпадает с количеством ключей")

// S3 представляет клиент для работы с AWS S3 хранилищем
type S3 struct {
	Client *s3.Client
	Bucket string
}

// NewS3 создает новый экземпляр S3 клиента
// ctx - контекст
// bucket - имя корзины S3
func NewS3(ctx context.Context, bucket string) (*S3, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &S3{
		Client: s3.NewFromConfig(cfg),
		Bucket: bucket,
	}, nil
}

// UploadFile загружает файл в S3
// ctx - контекст
// file - читатель файла
// key - ключ файла
// bucket - имя корзины S3
func (s *S3) UploadFile(ctx context.Context, file io.Reader, key string) error {
	harsher := sha256.New()
	sha256Hash := hex.EncodeToString(harsher.Sum(nil)) // Конвертируем хеш в строку
	contentType := getContentType(key)

	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:         aws.String(s.Bucket),
		Key:            aws.String(key),
		Body:           file,
		ContentType:    aws.String(contentType),
		ChecksumSHA256: aws.String(sha256Hash),
	})

	return err
}

// DownloadFile загружает файл из S3
// key - ключ файла
func (s *S3) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// DeleteFile удаляет файл из S3
// key - ключ файла
func (s *S3) DeleteFile(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	return err
}

// DeleteFiles удаляет несколько файлов из S3
// keys - ключи файлов
func (s *S3) DeleteFiles(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	objects := make([]types.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		objects = append(objects, types.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(s.Bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	}

	_, err := s.Client.DeleteObjects(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// UploadFiles загружает несколько файлов в S3
// files - читатели файлов
// keys - ключи файлов
func (s *S3) UploadFiles(ctx context.Context, files []io.Reader, keys []string) error {
	if len(files) != len(keys) {
		return ErrorFilesCountMismatch
	}

	wg := sync.WaitGroup{}
	for i, file := range files {
		wg.Add(1)
		go func(file io.Reader, key string) {
			defer wg.Done()
			if err := s.UploadFile(ctx, file, key); err != nil {
				fmt.Println(err)
			}
		}(file, keys[i])
	}
	wg.Wait()

	return nil
}

func getContentType(filename string) string {
	return mime.TypeByExtension(filepath.Ext(filename))
}
