package minio

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg config.MinioConfig) (*minio.Client, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioUser, cfg.MinioPassword, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateBucket(mc *minio.Client, cfg config.MinioConfig) error {
	ctx := context.Background()
	exists, err := mc.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return err
	}

	if !exists {
		err := mc.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
