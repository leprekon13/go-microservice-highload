package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IntegrationService struct {
	client     *minio.Client
	bucketName string
}

func NewIntegrationService(endpoint, accessKey, secretKey, bucket string) (*IntegrationService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &IntegrationService{
		client:     client,
		bucketName: bucket,
	}, nil
}

func (s *IntegrationService) InitBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *IntegrationService) ExportUsers(ctx context.Context, users any) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	objectName := fmt.Sprintf("users-export-%d.json", time.Now().Unix())
	reader := bytes.NewReader(data)

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/json",
	})
	return err
}
