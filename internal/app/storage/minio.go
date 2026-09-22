package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient() (*minio.Client, error) {

	client, err := minio.New(
		"localhost:9000",
		&minio.Options{
			Creds: credentials.NewStaticV4(
				"root",
				"rootpassword",
				"",
			),
			Secure: false,
		},
	)

	if err != nil {
		return nil, err
	}

	_, err = client.ListBuckets(
		context.Background(),
	)

	return client, err
}
