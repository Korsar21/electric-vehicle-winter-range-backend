package handler

import (
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/storage"
)

type APIHandler struct {
	Repository   *repository.Repository
	AssetStorage *storage.AssetStorage
}

func NewAPIHandler(
	repo *repository.Repository,
	assetStorage *storage.AssetStorage,
) *APIHandler {
	return &APIHandler{
		Repository:   repo,
		AssetStorage: assetStorage,
	}
}
