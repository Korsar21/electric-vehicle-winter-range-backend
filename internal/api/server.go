package api

import (
	"log"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/handler"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() error {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
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

	router.GET(
		"/vehicle-auxiliary-loads/feed",
		h.GetVehicleAuxiliaryLoadFeed,
	)

	router.GET(
		"/vehicle-auxiliary-loads/draft",
		h.GetDraftVehicleAuxiliaryLoad,
	)

	router.GET(
		"/vehicle-auxiliary-loads",
		h.GetVehicleAuxiliaryLoadGrid,
	)

	return router.Run(":8080")
}
