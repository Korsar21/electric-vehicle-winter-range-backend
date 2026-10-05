package api

import (
	"log"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/database"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/handler"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
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

	h := handler.NewHandler(repo)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static(
		"/static",
		"./resources",
	)

	// Лента
	router.GET(
		"/electric-car-loads/feed",
		h.GetElectricCarLoadFeed,
	)

	// Добавление / черновик
	router.GET(
		"/electric-car-loads/add",
		h.GetDraftElectricCarLoad,
	)

	// Создание черновика
	router.POST(
		"/electric-car-loads/add",
		h.CreateElectricCarLoad,
	)

	// Публикация черновика
	router.POST(
		"/electric-car-loads/:id/publish",
		h.PublishElectricCarLoad,
	)

	// Логическое удаление через SQL UPDATE
	router.POST(
		"/electric-car-loads/:id/delete",
		h.DeleteElectricCarLoad,
	)

	// Плитка + фильтрация по мощности
	router.GET(
		"/electric-car-loads",
		h.GetElectricCarLoadGrid,
	)

	return router.Run(":8080")
}
