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

	minioClient, err := storage.NewMinioClient()

	if err != nil {
		logrus.Error("minio error: ", err)
		return err
	}

	repo, err := repository.NewRepository(database.DB)
	if err != nil {
		logrus.Error("repository initialization error: ", err)
		return err
	}

	h := handler.NewHandler(
		repo,
		minioClient,
	)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static(
		"/static",
		"./resources",
	)

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

	router.GET(
		"/electric-car-loads/edit/:id",
		h.GetEditElectricCarLoad,
	)

	router.POST(
		"/electric-car-loads/edit/:id",
		h.UpdateElectricCarLoad,
	)

	router.POST(
		"/electric-car-loads/:id/delete",
		h.DeleteElectricCarLoad,
	)

	router.POST(
		"/electric-car-loads/:id/like",
		h.LikeElectricCarLoad,
	)

	router.GET(
		"/electric-car-loads/:id/edit",
		h.GetEditElectricCarLoad,
	)

	router.POST(
		"/electric-car-loads/:id/edit",
		h.UpdateElectricCarLoad,
	)

	router.GET(
		"/electric-car-loads",
		h.GetElectricCarLoadGrid,
	)
	return router.Run(":8080")
}
