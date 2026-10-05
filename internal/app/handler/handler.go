package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	defaultMaxPowerW = 10000
	maxAllowedPowerW = 10000

	//  2 авторизация не реализуется,
	// поэтому используем пользователя,
	// созданного Seed().
	defaultUserID uint = 1
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

	ImageURL string
	VideoURL string

	HasNext bool
}

func NewHandler(
	repo *repository.Repository,
) *Handler {

	return &Handler{
		Repository: repo,
	}
}

// makeMediaURL возвращает URL файла.
//
//	2 дефолтные файлы находятся непосредственно
//
// на SSR-сервере в resources/media/.
func makeMediaURL(key string, mediaType string) string {

	if key == "" {
		if mediaType == "video" {
			return "/static/media/default.mp4"
		}

		return "/static/media/default.jpg"
	}

	return "http://localhost:9000/ev-winter-assets/" + key
}

func makeElectricCarLoadView(
	load model.ElectricCarLoad,
) ElectricCarLoadView {

	return ElectricCarLoadView{
		ID:          int(load.ID),
		Name:        load.Name,
		Description: load.Description,
		Category:    load.Category,
		PowerDrawW:  load.PowerDrawW,
		Priority:    load.Priority,
		Status:      repository.ElectricCarLoadStatus(load.Status),
		LikeCount:   load.LikeCount,

		ImageURL: makeMediaURL(load.ImageKey, "image"),
		VideoURL: makeMediaURL(load.VideoKey, "video"),
	}
}

// GET /electric-car-loads/feed
func (h *Handler) GetElectricCarLoadFeed(ctx *gin.Context) {

	idQuery := ctx.Query("id")
	nextQuery := ctx.Query("next")

	var (
		load model.ElectricCarLoad
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

			load, err =
				h.Repository.GetNextPublishedElectricCarLoad(id)

		} else {

			load, err =
				h.Repository.GetPublishedElectricCarLoadByID(id)
		}
	}

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusNotFound,
			"electric car load not found",
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

// GET /electric-car-loads/add
func (h *Handler) GetDraftElectricCarLoad(ctx *gin.Context) {

	load, err :=
		h.Repository.GetDraftElectricCarLoad(defaultUserID)

	// сли черновика нет — это нормальная ситуация.
	// Страница просто открывается пустой.
	if errors.Is(err, gorm.ErrRecordNotFound) {

		ctx.HTML(
			http.StatusOK,
			"electric_car_loads_add.html",
			gin.H{
				"Load": ElectricCarLoadView{
					Status:   repository.StatusDraft,
					ImageURL: "/static/media/default.jpg",
					VideoURL: "/static/media/default.mp4",
				},
				"HasDraft": false,
			},
		)

		return
	}

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"failed to load draft",
		)

		return
	}

	view := makeElectricCarLoadView(load)

	ctx.HTML(
		http.StatusOK,
		"electric_car_loads_add.html",
		gin.H{
			"Load":     view,
			"HasDraft": true,
		},
	)
}

// GET /electric-car-loads
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
				"maximum power draw must be between 0 and 10000 W",
			)

			return
		}

		maxPowerW = value
	}

	loads, err :=
		h.Repository.GetPublishedElectricCarLoads(maxPowerW)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"failed to load electric car loads",
		)

		return
	}

	views :=
		make([]ElectricCarLoadView, 0, len(loads))

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

// POST /electric-car-loads/add
//
// а первом этапе пользователь указывает только название.
// ото и видео в 2  передаются и  сохраняются.
// нопка "алее" создаёт черновик.
func (h *Handler) CreateElectricCarLoad(ctx *gin.Context) {

	name := ctx.PostForm("name")

	if name == "" {

		ctx.String(
			http.StatusBadRequest,
			"name is required",
		)

		return
	}

	load := model.ElectricCarLoad{
		Name:      name,
		Status:    string(repository.StatusDraft),
		CreatorID: defaultUserID,
	}

	err :=
		h.Repository.CreateElectricCarLoad(&load)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"failed to create draft",
		)

		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electric-car-loads/add",
	)
}

// POST /electric-car-loads/:id/publish
//
// аполняются поля предметной области,
// после чего черновик переводится в published.
func (h *Handler) PublishElectricCarLoad(ctx *gin.Context) {

	id, err :=
		strconv.Atoi(ctx.Param("id"))

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"invalid id",
		)

		return
	}

	power, err :=
		strconv.Atoi(ctx.PostForm("power"))

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"invalid power value",
		)

		return
	}

	err =
		h.Repository.PublishElectricCarLoad(
			id,
			ctx.PostForm("description"),
			ctx.PostForm("category"),
			power,
			ctx.PostForm("priority"),
		)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"failed to publish electric car load",
		)

		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electric-car-loads",
	)
}

// POST /electric-car-loads/:id/delete
//
// огическое удаление через SQL UPDATE.
func (h *Handler) DeleteElectricCarLoad(ctx *gin.Context) {

	id, err :=
		strconv.Atoi(ctx.Param("id"))

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"invalid id",
		)

		return
	}

	err =
		h.Repository.DeleteElectricCarLoad(id)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"failed to delete electric car load",
		)

		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electric-car-loads",
	)
}
