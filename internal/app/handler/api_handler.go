package handler

import (
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/auth"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/storage"
)

type APIHandler struct {
	Repository   *repository.Repository
	AssetStorage *storage.AssetStorage
	Auth         *auth.Manager
}

func NewAPIHandler(
	repo *repository.Repository,
	assetStorage *storage.AssetStorage,
	authManager *auth.Manager,
) *APIHandler {
	return &APIHandler{
		Repository:   repo,
		AssetStorage: assetStorage,
		Auth:         authManager,
	}
}
