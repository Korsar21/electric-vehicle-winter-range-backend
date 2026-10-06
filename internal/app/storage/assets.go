package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const assetsBucket = "ev-winter-assets"

type AssetStorage struct {
	client *minio.Client
}

func NewAssetStorage(client *minio.Client) *AssetStorage {
	return &AssetStorage{
		client: client,
	}
}

func (s *AssetStorage) Save(
	ctx context.Context,
	reader io.Reader,
	size int64,
	contentType string,
	originalName string,
) (string, error) {
	ext := filepath.Ext(originalName)

	key := fmt.Sprintf(
		"%s%s",
		uuid.New().String(),
		ext,
	)

	_, err := s.client.PutObject(
		ctx,
		assetsBucket,
		key,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *AssetStorage) Delete(
	ctx context.Context,
	key string,
) error {
	if key == "" {
		return nil
	}

	return s.client.RemoveObject(
		ctx,
		assetsBucket,
		key,
		minio.RemoveObjectOptions{},
	)
}
