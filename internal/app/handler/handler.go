package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	minioBaseURL     = "http://localhost:9000"
	minioBucketName  = "ev-winter-assets"
	defaultMaxPowerW = 5000
	maxAllowedPowerW = 5000
)

type Handler struct {
	Repository *repository.Repository
}

type VehicleAuxiliaryLoadView struct {
	ID          int
	Name        string
	Description string
	Category    string
	PowerDrawW  int
	Status      repository.AuxiliaryLoadStatus
	LikeCount   int
	ImageURL    string
	VideoURL    string
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func makeMediaURL(key string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		minioBaseURL,
		minioBucketName,
		key,
	)
}

func makeVehicleAuxiliaryLoadView(
	load repository.VehicleAuxiliaryLoad,
) VehicleAuxiliaryLoadView {
	return VehicleAuxiliaryLoadView{
		ID:          load.ID,
		Name:        load.Name,
		Description: load.Description,
		Category:    load.Category,
		PowerDrawW:  load.PowerDrawW,
		Status:      load.Status,

		LikeCount: len(load.LikeUserIDs),

		ImageURL: makeMediaURL(load.ImageKey),
		VideoURL: makeMediaURL(load.VideoKey),
	}
}

func (h *Handler) GetVehicleAuxiliaryLoadFeed(ctx *gin.Context) {
	idQuery := ctx.Query("id")
	nextQuery := ctx.Query("next")

	var (
		load repository.VehicleAuxiliaryLoad
		err  error
	)

	if idQuery == "" {
		load, err = h.Repository.GetFirstPublishedVehicleAuxiliaryLoad()
	} else {
		id, parseErr := strconv.Atoi(idQuery)
		if parseErr != nil {
			ctx.String(
				http.StatusBadRequest,
				"invalid vehicle auxiliary load id",
			)
			return
		}

		if nextQuery == "true" {
			load, err = h.Repository.GetNextPublishedVehicleAuxiliaryLoad(id)
		} else {
			load, err = h.Repository.GetPublishedVehicleAuxiliaryLoadByID(id)
		}
	}

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	view := makeVehicleAuxiliaryLoadView(load)

	ctx.HTML(
		http.StatusOK,
		"feed.html",
		gin.H{
			"Load": view,
		},
	)
}

func (h *Handler) GetDraftVehicleAuxiliaryLoad(ctx *gin.Context) {
	load, err := h.Repository.GetDraftVehicleAuxiliaryLoad()
	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	view := makeVehicleAuxiliaryLoadView(load)

	ctx.HTML(
		http.StatusOK,
		"add.html",
		gin.H{
			"Load": view,
		},
	)
}

func (h *Handler) GetVehicleAuxiliaryLoadGrid(ctx *gin.Context) {
	maxPowerW := defaultMaxPowerW

	maxPowerQuery := ctx.Query("maxPowerW")

	if maxPowerQuery != "" {
		value, err := strconv.Atoi(maxPowerQuery)
		if err != nil {
			ctx.String(
				http.StatusBadRequest,
				"invalid maximum power draw",
			)
			return
		}

		if value < 0 || value > maxAllowedPowerW {
			ctx.String(
				http.StatusBadRequest,
				"maximum power draw must be between 0 and 5000 W",
			)
			return
		}

		maxPowerW = value
	}

	loads, err := h.Repository.GetPublishedVehicleAuxiliaryLoads(maxPowerW)
	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"failed to load vehicle auxiliary loads",
		)
		return
	}

	views := make([]VehicleAuxiliaryLoadView, 0, len(loads))

	for _, load := range loads {
		views = append(
			views,
			makeVehicleAuxiliaryLoadView(load),
		)
	}

	ctx.HTML(
		http.StatusOK,
		"grid.html",
		gin.H{
			"Loads":     views,
			"MaxPowerW": maxPowerW,
		},
	)
}
