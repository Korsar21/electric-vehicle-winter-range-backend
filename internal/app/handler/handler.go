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

type ElectricCarLoadView struct {
	ID          int
	Name        string
	Description string
	Category    string
	PowerDrawW  int
	Priority    string
	Status      repository.ElectricCarLoadStatus
	LikeCount   int

	ImageKey string
	VideoKey string

	ImageURL string
	VideoURL string

	HasNext bool
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

func makeElectricCarLoadView(
	load repository.ElectricCarLoad,
) ElectricCarLoadView {
	return ElectricCarLoadView{
		ID:          load.ID,
		Name:        load.Name,
		Description: load.Description,
		Category:    load.Category,
		PowerDrawW:  load.PowerDrawW,
		Priority:    load.Priority,
		Status:      load.Status,

		LikeCount: len(load.LikeUserIDs),

		ImageKey: load.ImageKey,
		VideoKey: load.VideoKey,

		ImageURL: makeMediaURL(load.ImageKey),
		VideoURL: makeMediaURL(load.VideoKey),
	}
}

func (h *Handler) GetElectricCarLoadFeed(ctx *gin.Context) {
	idQuery := ctx.Query("id")
	nextQuery := ctx.Query("next")

	var (
		load repository.ElectricCarLoad
		err  error
	)

	if idQuery == "" {
		load, err = h.Repository.GetFirstPublishedElectricCarLoad()
	} else {
		id, parseErr := strconv.Atoi(idQuery)
		if parseErr != nil {
			ctx.String(
				http.StatusBadRequest,
				"invalid electric car load id",
			)
			return
		}

		if nextQuery == "true" {
			load, err = h.Repository.GetNextPublishedElectricCarLoad(id)
		} else {
			load, err = h.Repository.GetPublishedElectricCarLoadByID(id)
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

	view := makeElectricCarLoadView(load)
	view.HasNext = true

	ctx.HTML(
		http.StatusOK,
		"electric_car_loads_feed.html",
		gin.H{
			"Load": view,
		},
	)
}

func (h *Handler) GetDraftElectricCarLoad(ctx *gin.Context) {
	load, err := h.Repository.GetDraftElectricCarLoad()
	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	view := makeElectricCarLoadView(load)

	ctx.HTML(
		http.StatusOK,
		"electric_car_loads_add.html",
		gin.H{
			"Load": view,
		},
	)
}

func (h *Handler) GetElectricCarLoadGrid(ctx *gin.Context) {
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

	loads, err := h.Repository.GetPublishedElectricCarLoads(maxPowerW)
	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"failed to load electric car loads",
		)
		return
	}

	views := make([]ElectricCarLoadView, 0, len(loads))

	for _, load := range loads {
		views = append(
			views,
			makeElectricCarLoadView(load),
		)
	}

	ctx.HTML(
		http.StatusOK,
		"electric_car_loads_grid.html",
		gin.H{
			"Loads":     views,
			"MaxPowerW": maxPowerW,
		},
	)
}
