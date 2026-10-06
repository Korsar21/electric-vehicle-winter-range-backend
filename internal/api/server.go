package api

import (
	"log"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/database"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/handler"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() error {

	log.Println("Starting server")

	database.Connect()

	database.Migrate()

	database.Seed()

	repo, err := repository.NewRepository(database.DB)

	if err != nil {
		logrus.Error("repository initialization error: ", err)
		return err
	}

	minioClient, err := storage.NewMinioClient()

	if err != nil {
		logrus.Error("MinIO initialization error: ", err)
		return err
	}

	assetStorage := storage.NewAssetStorage(minioClient)

	h := handler.NewHandler(repo)

	apiHandler := handler.NewAPIHandler(
		repo,
		assetStorage,
	)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static(
		"/static",
		"./resources",
	)

	// =========================
	// SSR / HTML — 2
	// =========================

	router.GET(
		"/electric-car-loads/feed",
		h.GetElectricCarLoadFeed,
	)

	router.GET(
		"/electric-car-loads/add",
		h.GetDraftElectricCarLoad,
	)

	router.POST(
		"/electric-car-loads/add",
		h.CreateElectricCarLoad,
	)

	router.POST(
		"/electric-car-loads/:id/publish",
		h.PublishElectricCarLoad,
	)

	router.POST(
		"/electric-car-loads/:id/delete",
		h.DeleteElectricCarLoad,
	)

	router.GET(
		"/electric-car-loads",
		h.GetElectricCarLoadGrid,
	)

	// =========================
	// REST API — 3
	// =========================

	// Список опубликованных услуг + фильтрация.
	router.GET(
		"/api/electric-car-loads",
		apiHandler.GetElectricCarLoads,
	)

	// ента:
	// /api/electric-car-loads/feed
	// /api/electric-car-loads/feed?id=1
	// /api/electric-car-loads/feed?id=1&next=true
	router.GET(
		"/api/electric-car-loads/feed",
		apiHandler.GetElectricCarLoadFeed,
	)

	// ерновик текущего пользователя.
	router.GET(
		"/api/electric-car-loads/draft",
		apiHandler.GetDraftElectricCarLoadAPI,
	)

	// Создание черновика + загрузка image/video.
	router.POST(
		"/api/electric-car-loads",
		apiHandler.CreateElectricCarLoadAPI,
	)

	// убликация черновика.
	router.PUT(
		"/api/electric-car-loads/:id",
		apiHandler.PublishElectricCarLoadAPI,
	)

	// Soft delete собственной услуги.
	router.DELETE(
		"/api/electric-car-loads/:id",
		apiHandler.DeleteElectricCarLoadAPI,
	)

	// становка/отмена лайка.
	router.POST(
		"/api/electric-car-loads/:id/like",
		apiHandler.LikeElectricCarLoadAPI,
	)
	// =========================
	// USER API — 3
	// =========================

	router.POST(
		"/api/users/register",
		apiHandler.RegisterUser,
	)

	router.POST(
		"/api/users/auth",
		apiHandler.AuthenticateUser,
	)

	router.POST(
		"/api/users/logout",
		apiHandler.LogoutUser,
	)

	return router.Run(":8080")
}
