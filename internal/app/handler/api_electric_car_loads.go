package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/auth"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/serializer"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	apiDefaultMaxPowerW = 10000
	apiMaxAllowedPowerW = 10000
)

func apiError(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"status":  "fail",
		"message": message,
	})
}

func makeElectricCarLoadResponse(
	h *APIHandler,
	load model.ElectricCarLoad,
) serializer.ElectricCarLoadResponse {
	currentUserID := auth.CurrentUserID()

	isCreator := 0
	if load.CreatorID == currentUserID {
		isCreator = 1
	}

	isLiked := 0

	liked, err := h.Repository.IsLiked(
		currentUserID,
		load.ID,
	)

	if err == nil && liked {
		isLiked = 1
	}

	return serializer.ElectricCarLoadResponse{
		ID:          load.ID,
		Name:        load.Name,
		Description: load.Description,
		Category:    load.Category,
		PowerDrawW:  load.PowerDrawW,
		Priority:    load.Priority,
		Status:      load.Status,
		LikeCount:   load.LikeCount,
		ImageURL:    makeMediaURL(load.ImageKey, "image"),
		VideoURL:    makeMediaURL(load.VideoKey, "video"),
		IsCreator:   isCreator,
		IsLiked:     isLiked,
	}
}

// GET /api/electric-car-loads
//
// Только опубликованные услуги.
// оддерживается фильтрация по максимальной мощности.
func (h *APIHandler) GetElectricCarLoads(ctx *gin.Context) {
	maxPowerW := apiDefaultMaxPowerW

	maxPowerQuery := ctx.Query("maxPowerW")

	if maxPowerQuery != "" {
		value, err := strconv.Atoi(maxPowerQuery)

		if err != nil {
			apiError(
				ctx,
				http.StatusBadRequest,
				"invalid maximum power draw",
			)
			return
		}

		if value < 0 || value > apiMaxAllowedPowerW {
			apiError(
				ctx,
				http.StatusBadRequest,
				"maximum power draw must be between 0 and 10000 W",
			)
			return
		}

		maxPowerW = value
	}

	loads, err := h.Repository.GetPublishedElectricCarLoads(maxPowerW)

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to load electric car loads",
		)
		return
	}

	response := make(
		[]serializer.ElectricCarLoadResponse,
		0,
		len(loads),
	)

	for _, load := range loads {
		response = append(
			response,
			makeElectricCarLoadResponse(h, load),
		)
	}

	ctx.JSON(http.StatusOK, response)
}

// GET /api/electric-car-loads/feed
//
// ез id — первый опубликованный объект.
// С id — конкретный объект.
// С id и next=true — следующий опубликованный объект.
func (h *APIHandler) GetElectricCarLoadFeed(ctx *gin.Context) {
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

		if parseErr != nil || id <= 0 {
			apiError(
				ctx,
				http.StatusBadRequest,
				"invalid id",
			)
			return
		}

		if nextQuery == "true" {
			load, err = h.Repository.GetNextPublishedElectricCarLoad(id)
		} else {
			load, err = h.Repository.GetPublishedElectricCarLoadByID(id)
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		apiError(
			ctx,
			http.StatusNotFound,
			"electric car load not found",
		)
		return
	}

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to load electric car load",
		)
		return
	}

	response := makeElectricCarLoadResponse(h, load)

	// роверяем, существует ли следующий опубликованный объект.
	nextLoad, nextErr := h.Repository.GetNextPublishedElectricCarLoad(
		int(load.ID),
	)

	if nextErr == nil && nextLoad.ID != load.ID {
		response.HasNext = true
	}

	ctx.JSON(http.StatusOK, response)
}

// GET /api/electric-car-loads/draft
//
// озвращает черновик текущего пользователя.
//
//	пользователя может быть не более одного черновика.
func (h *APIHandler) GetDraftElectricCarLoadAPI(ctx *gin.Context) {
	load, err := h.Repository.GetDraftElectricCarLoad(
		auth.CurrentUserID(),
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		apiError(
			ctx,
			http.StatusNotFound,
			"draft electric car load not found",
		)
		return
	}

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to load draft",
		)
		return
	}

	response := makeElectricCarLoadResponse(h, load)

	ctx.JSON(http.StatusOK, response)
}

// POST /api/electric-car-loads
//
// Создаёт черновик текущего пользователя.
// name передаётся как текстовое поле.
// image и video передаются как реальные multipart-файлы.
func (h *APIHandler) CreateElectricCarLoadAPI(ctx *gin.Context) {
	currentUserID := auth.CurrentUserID()

	_, err := h.Repository.GetDraftElectricCarLoad(currentUserID)

	if err == nil {
		apiError(
			ctx,
			http.StatusConflict,
			"current user already has a draft",
		)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to check existing draft",
		)
		return
	}

	name := ctx.PostForm("name")

	if name == "" {
		apiError(
			ctx,
			http.StatusBadRequest,
			"name is required",
		)
		return
	}

	imageHeader, imageErr := ctx.FormFile("image")

	if imageErr != nil {
		apiError(
			ctx,
			http.StatusBadRequest,
			"image file is required",
		)
		return
	}

	videoHeader, videoErr := ctx.FormFile("video")

	if videoErr != nil {
		apiError(
			ctx,
			http.StatusBadRequest,
			"video file is required",
		)
		return
	}

	imageContentType := imageHeader.Header.Get("Content-Type")
	videoContentType := videoHeader.Header.Get("Content-Type")

	if len(imageContentType) < 6 ||
		imageContentType[:6] != "image/" {
		apiError(
			ctx,
			http.StatusBadRequest,
			"image must be an image file",
		)
		return
	}

	if len(videoContentType) < 6 ||
		videoContentType[:6] != "video/" {
		apiError(
			ctx,
			http.StatusBadRequest,
			"video must be a video file",
		)
		return
	}

	imageFile, err := imageHeader.Open()

	if err != nil {
		apiError(
			ctx,
			http.StatusBadRequest,
			"failed to open image",
		)
		return
	}
	defer imageFile.Close()

	imageKey, err := h.AssetStorage.Save(
		ctx.Request.Context(),
		imageFile,
		imageHeader.Size,
		imageContentType,
		imageHeader.Filename,
	)

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to save image",
		)
		return
	}

	videoFile, err := videoHeader.Open()

	if err != nil {
		_ = h.AssetStorage.Delete(
			ctx.Request.Context(),
			imageKey,
		)

		apiError(
			ctx,
			http.StatusBadRequest,
			"failed to open video",
		)
		return
	}
	defer videoFile.Close()

	videoKey, err := h.AssetStorage.Save(
		ctx.Request.Context(),
		videoFile,
		videoHeader.Size,
		videoContentType,
		videoHeader.Filename,
	)

	if err != nil {
		_ = h.AssetStorage.Delete(
			ctx.Request.Context(),
			imageKey,
		)

		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to save video",
		)
		return
	}

	load := model.ElectricCarLoad{
		Name:      name,
		Status:    string(repository.StatusDraft),
		ImageKey:  imageKey,
		VideoKey:  videoKey,
		CreatorID: currentUserID,
	}

	if err := h.Repository.CreateElectricCarLoad(&load); err != nil {
		_ = h.AssetStorage.Delete(
			ctx.Request.Context(),
			imageKey,
		)

		_ = h.AssetStorage.Delete(
			ctx.Request.Context(),
			videoKey,
		)

		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to create electric car load",
		)
		return
	}

	response := makeElectricCarLoadResponse(h, load)

	ctx.JSON(
		http.StatusCreated,
		response,
	)
}

// PUT /api/electric-car-loads/:id
//
// убликация черновика.
// Только создатель может опубликовать собственный черновик.
func (h *APIHandler) PublishElectricCarLoadAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		apiError(
			ctx,
			http.StatusBadRequest,
			"invalid id",
		)
		return
	}

	var request serializer.ElectricCarLoadPublishRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		apiError(
			ctx,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if request.PowerDrawW < 0 ||
		request.PowerDrawW > apiMaxAllowedPowerW {
		apiError(
			ctx,
			http.StatusBadRequest,
			"power_draw_w must be between 0 and 10000 W",
		)
		return
	}

	err = h.Repository.PublishElectricCarLoad(
		id,
		auth.CurrentUserID(),
		request.Description,
		request.Category,
		request.PowerDrawW,
		request.Priority,
	)

	if errors.Is(err, repository.ErrElectricCarLoadNotFound) {
		apiError(
			ctx,
			http.StatusNotFound,
			"draft electric car load not found",
		)
		return
	}

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to publish electric car load",
		)
		return
	}

	load, err := h.Repository.GetPublishedElectricCarLoadByID(id)

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to load published electric car load",
		)
		return
	}

	response := makeElectricCarLoadResponse(h, load)

	ctx.JSON(http.StatusOK, response)
}

// DELETE /api/electric-car-loads/:id
//
// Только soft delete.
// далять можно только собственную услугу.
func (h *APIHandler) DeleteElectricCarLoadAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		apiError(
			ctx,
			http.StatusBadRequest,
			"invalid id",
		)
		return
	}

	err = h.Repository.DeleteElectricCarLoad(
		id,
		auth.CurrentUserID(),
	)

	if errors.Is(err, repository.ErrElectricCarLoadNotFound) {
		apiError(
			ctx,
			http.StatusNotFound,
			"electric car load not found",
		)
		return
	}

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to delete electric car load",
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "electric car load deleted",
	})
}

// POST /api/electric-car-loads/:id/like
//
// like=1 — поставить лайк.
// like=0 — отменить лайк.
func (h *APIHandler) LikeElectricCarLoadAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		apiError(
			ctx,
			http.StatusBadRequest,
			"invalid id",
		)
		return
	}

	var request serializer.LikeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		apiError(
			ctx,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if request.Like != 0 && request.Like != 1 {
		apiError(
			ctx,
			http.StatusBadRequest,
			"like must be 0 or 1",
		)
		return
	}

	load, err := h.Repository.GetPublishedElectricCarLoadByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		apiError(
			ctx,
			http.StatusNotFound,
			"electric car load not found",
		)
		return
	}

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to load electric car load",
		)
		return
	}

	err = h.Repository.SetLike(
		auth.CurrentUserID(),
		load.ID,
		request.Like,
	)

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to update like",
		)
		return
	}

	liked, err := h.Repository.IsLiked(
		auth.CurrentUserID(),
		load.ID,
	)

	if err != nil {
		apiError(
			ctx,
			http.StatusInternalServerError,
			"failed to check like",
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"is_liked": boolToInt(liked),
	})
}

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}
