package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

const (
	minioBaseURL     = "http://localhost:9000"
	minioBucketName  = "ev-winter-assets"
	defaultMaxPowerW = 10000
	maxAllowedPowerW = 10000
)

type Handler struct {
	Repository *repository.Repository
	Minio      *minio.Client
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

func NewHandler(
	r *repository.Repository,
	minioClient *minio.Client,
) *Handler {

	return &Handler{
		Repository: r,
		Minio:      minioClient,
	}
}

func makeMediaURL(key string) string {

	if key == "" {
		key = "default.jpg"
	}

	return fmt.Sprintf(
		"%s/%s/%s",
		minioBaseURL,
		minioBucketName,
		key,
	)
}

func (h *Handler) uploadFile(
	fileHeader *multipart.FileHeader,
) error {

	file, err := fileHeader.Open()

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = h.Minio.PutObject(
		context.Background(),
		minioBucketName,
		fileHeader.Filename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{},
	)

	return err
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

		LikeCount: load.LikeCount,

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

func (h *Handler) CreateElectricCarLoad(ctx *gin.Context) {

	imageKey := ""
	videoKey := ""

	image, err := ctx.FormFile("image")

	if err == nil {

		err = h.uploadFile(image)

		if err != nil {

			ctx.String(
				http.StatusInternalServerError,
				err.Error(),
			)

			return
		}

		imageKey = image.Filename
	}

	video, err := ctx.FormFile("video")

	if err == nil {

		err = h.uploadFile(video)

		if err != nil {

			ctx.String(
				http.StatusInternalServerError,
				err.Error(),
			)

			return
		}

		videoKey = video.Filename
	}

	power, err := strconv.Atoi(
		ctx.PostForm("power"),
	)

	if err != nil {
		ctx.String(
			http.StatusBadRequest,
			"invalid power value",
		)
		return
	}

	load := repository.ElectricCarLoad{

		Name: ctx.PostForm("name"),

		Description: ctx.PostForm("description"),

		Category: ctx.PostForm("category"),

		PowerDrawW: power,

		Priority: ctx.PostForm("priority"),

		Status: repository.StatusPublished,

		ImageKey: imageKey,

		VideoKey: videoKey,

		LikeCount: 0,
	}

	err = h.Repository.CreateElectricCarLoad(load)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electric-car-loads",
	)
}

func (h *Handler) DeleteElectricCarLoad(ctx *gin.Context) {

	id, err := strconv.Atoi(
		ctx.Param("id"),
	)

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"invalid id",
		)

		return
	}

	err = h.Repository.DeleteElectricCarLoad(id)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electric-car-loads",
	)
}

func (h *Handler) GetEditElectricCarLoad(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"invalid id",
		)

		return
	}

	load, err := h.Repository.GetElectricCarLoadByID(id)

	if err != nil {

		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)

		return
	}

	view := makeElectricCarLoadView(load)

	ctx.HTML(
		http.StatusOK,
		"electric_car_loads_edit.html",
		gin.H{
			"Load": view,
		},
	)

}

func (h *Handler) UpdateElectricCarLoad(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"invalid id",
		)

		return
	}

	load, err := h.Repository.GetElectricCarLoadByID(id)

	if err != nil {

		ctx.String(
			http.StatusNotFound,
			err.Error(),
		)

		return
	}

	load.Name = ctx.PostForm("name")

	load.Description = ctx.PostForm("description")

	load.Category = ctx.PostForm("category")

	load.Priority = ctx.PostForm("priority")

	power, err := strconv.Atoi(
		ctx.PostForm("power"),
	)

	if err == nil {
		load.PowerDrawW = power
	}

	// обновление изображения
	image, err := ctx.FormFile("image")

	if err == nil {

		err = h.uploadFile(image)

		if err != nil {

			ctx.String(
				http.StatusInternalServerError,
				err.Error(),
			)

			return
		}

		load.ImageKey = image.Filename
	}

	// обновление видео
	video, err := ctx.FormFile("video")

	if err == nil {

		err = h.uploadFile(video)

		if err != nil {

			ctx.String(
				http.StatusInternalServerError,
				err.Error(),
			)

			return
		}

		load.VideoKey = video.Filename
	}

	err = h.Repository.UpdateElectricCarLoad(load)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electric-car-loads",
	)

}

func (h *Handler) LikeElectricCarLoad(ctx *gin.Context) {

	id, err := strconv.Atoi(
		ctx.Param("id"),
	)

	if err != nil {

		ctx.String(
			http.StatusBadRequest,
			"invalid id",
		)

		return
	}

	err = h.Repository.LikeElectricCarLoad(id)

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electric-car-loads/feed?id="+strconv.Itoa(id),
	)
}
