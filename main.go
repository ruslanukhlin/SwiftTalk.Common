package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

    // Подгружаем конфигурацию из ~/.aws/*
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

    // Создаем клиента для доступа к хранилищу S3
    client := s3.NewFromConfig(cfg)

    // Запрашиваем список бакетов
    result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
    if err != nil {
        log.Fatal(err)
    }

    for _, bucket := range result.Buckets {
        log.Printf("bucket=%s creation time=%s", aws.ToString(bucket.Name), bucket.CreationDate.Local().Format("2006-01-02 15:04:05 Monday"))
    }
}	