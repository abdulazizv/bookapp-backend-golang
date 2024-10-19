package util

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.com/bookapp/config"
	"gitlab.com/bookapp/pkg/logger"
)

func Uploader(objectName, filePath string) (string, error) {
	ctx := context.Background()

	cfg := config.Load()

	minioClient, err := minio.New(cfg.MinioEnpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSsl,
	})
	if err != nil {
		logger.Error(err)
	}

	bucketName := "bookstore"
	location := "us-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
			return "", errBucketExists
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
		return "", err
	}

	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return objectName, nil
}
