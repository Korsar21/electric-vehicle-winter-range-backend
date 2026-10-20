package api

import (
	"log"

	_ "github.com/Korsar21/electric-vehicle-winter-range-backend/docs"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/auth"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/database"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/handler"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	authManager, err := auth.NewManager()
	if err != nil {
		logrus.Error("auth initialization error: ", err)
		return err
	}

	h := handler.NewHandler(repo)

	apiHandler := handler.NewAPIHandler(
		repo,
		assetStorage,
		authManager,
	)

	router := gin.Default()

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	router.LoadHTMLGlob("templates/*")

	router.Static(
		"/static",
		"./resources",
	)

	// SSR / HTML
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

	// Public REST API.
	// GET works both for guest and creator.
	router.GET(
		"/api/electric-car-loads",
		authManager.Optional(),
		apiHandler.GetElectricCarLoads,
	)

	router.GET(
		"/api/electric-car-loads/feed",
		authManager.Optional(),
		apiHandler.GetElectricCarLoadFeed,
	)

	router.POST(
		"/api/users/register",
		apiHandler.RegisterUser,
	)

	router.POST(
		"/api/users/auth",
		apiHandler.AuthenticateUser,
	)

	// Creator REST API.
	creator := router.Group("/api")
	creator.Use(authManager.Required())

	creator.GET(
		"/electric-car-loads/draft",
		apiHandler.GetDraftElectricCarLoadAPI,
	)

	creator.POST(
		"/electric-car-loads",
		apiHandler.CreateElectricCarLoadAPI,
	)

	creator.PUT(
		"/electric-car-loads/:id",
		apiHandler.PublishElectricCarLoadAPI,
	)

	creator.DELETE(
		"/electric-car-loads/:id",
		apiHandler.DeleteElectricCarLoadAPI,
	)

	creator.POST(
		"/electric-car-loads/:id/like",
		apiHandler.LikeElectricCarLoadAPI,
	)

	creator.POST(
		"/users/logout",
		apiHandler.LogoutUser,
	)

	return router.Run(":8080")
}
