package s3

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Client *s3.Client
	Bucket string
}

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

func (s *S3) UploadFile(ctx context.Context, file io.Reader, key string) error {
	harsher := sha256.New()
	sha256Hash := hex.EncodeToString(harsher.Sum(nil)) // Конвертируем хеш в строку
	contentType := getContentType(key)
	
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   file,		
		ContentType:    aws.String(contentType),
		ChecksumSHA256: aws.String(sha256Hash),
	})

	return err
}

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

func (s *S3) DeleteFile(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	return err
}

func getContentType(filename string) string {
	return mime.TypeByExtension(filepath.Ext(filename))
}
