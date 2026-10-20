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
	ctx *gin.Context,
	h *APIHandler,
	load model.ElectricCarLoad,
) serializer.ElectricCarLoadResponse {
	currentUserID, _ := auth.UserIDFromContext(ctx)

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
// РўРѕР»СЊРєРѕ РѕРїСѓР±Р»РёРєРѕРІР°РЅРЅС‹Рµ СѓСЃР»СѓРіРё.
// РѕРґРґРµСЂР¶РёРІР°РµС‚СЃСЏ С„РёР»СЊС‚СЂР°С†РёСЏ РїРѕ РјР°РєСЃРёРјР°Р»СЊРЅРѕР№ РјРѕС‰РЅРѕСЃС‚Рё.
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
			makeElectricCarLoadResponse(ctx, h, load),
		)
	}

	ctx.JSON(http.StatusOK, response)
}

// GET /api/electric-car-loads/feed
//
// РµР· id вЂ” РїРµСЂРІС‹Р№ РѕРїСѓР±Р»РёРєРѕРІР°РЅРЅС‹Р№ РѕР±СЉРµРєС‚.
// РЎ id вЂ” РєРѕРЅРєСЂРµС‚РЅС‹Р№ РѕР±СЉРµРєС‚.
// РЎ id Рё next=true вЂ” СЃР»РµРґСѓСЋС‰РёР№ РѕРїСѓР±Р»РёРєРѕРІР°РЅРЅС‹Р№ РѕР±СЉРµРєС‚.
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

	response := makeElectricCarLoadResponse(ctx, h, load)

	// СЂРѕРІРµСЂСЏРµРј, СЃСѓС‰РµСЃС‚РІСѓРµС‚ Р»Рё СЃР»РµРґСѓСЋС‰РёР№ РѕРїСѓР±Р»РёРєРѕРІР°РЅРЅС‹Р№ РѕР±СЉРµРєС‚.
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
// РѕР·РІСЂР°С‰Р°РµС‚ С‡РµСЂРЅРѕРІРёРє С‚РµРєСѓС‰РµРіРѕ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ.
//
//	РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ РјРѕР¶РµС‚ Р±С‹С‚СЊ РЅРµ Р±РѕР»РµРµ РѕРґРЅРѕРіРѕ С‡РµСЂРЅРѕРІРёРєР°.
func (h *APIHandler) GetDraftElectricCarLoadAPI(ctx *gin.Context) {
	load, err := h.Repository.GetDraftElectricCarLoad(
		auth.MustUserID(ctx),
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

	response := makeElectricCarLoadResponse(ctx, h, load)

	ctx.JSON(http.StatusOK, response)
}

// POST /api/electric-car-loads
//
// РЎРѕР·РґР°С‘С‚ С‡РµСЂРЅРѕРІРёРє С‚РµРєСѓС‰РµРіРѕ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ.
// name РїРµСЂРµРґР°С‘С‚СЃСЏ РєР°Рє С‚РµРєСЃС‚РѕРІРѕРµ РїРѕР»Рµ.
// image Рё video РїРµСЂРµРґР°СЋС‚СЃСЏ РєР°Рє СЂРµР°Р»СЊРЅС‹Рµ multipart-С„Р°Р№Р»С‹.
func (h *APIHandler) CreateElectricCarLoadAPI(ctx *gin.Context) {
	currentUserID, _ := auth.UserIDFromContext(ctx)

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

	response := makeElectricCarLoadResponse(ctx, h, load)

	ctx.JSON(
		http.StatusCreated,
		response,
	)
}

// PUT /api/electric-car-loads/:id
//
// СѓР±Р»РёРєР°С†РёСЏ С‡РµСЂРЅРѕРІРёРєР°.
// РўРѕР»СЊРєРѕ СЃРѕР·РґР°С‚РµР»СЊ РјРѕР¶РµС‚ РѕРїСѓР±Р»РёРєРѕРІР°С‚СЊ СЃРѕР±СЃС‚РІРµРЅРЅС‹Р№ С‡РµСЂРЅРѕРІРёРє.
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
		auth.MustUserID(ctx),
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

	response := makeElectricCarLoadResponse(ctx, h, load)

	ctx.JSON(http.StatusOK, response)
}

// DELETE /api/electric-car-loads/:id
//
// РўРѕР»СЊРєРѕ soft delete.
// РґР°Р»СЏС‚СЊ РјРѕР¶РЅРѕ С‚РѕР»СЊРєРѕ СЃРѕР±СЃС‚РІРµРЅРЅСѓСЋ СѓСЃР»СѓРіСѓ.
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
		auth.MustUserID(ctx),
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
// like=1 вЂ” РїРѕСЃС‚Р°РІРёС‚СЊ Р»Р°Р№Рє.
// like=0 вЂ” РѕС‚РјРµРЅРёС‚СЊ Р»Р°Р№Рє.
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
		auth.MustUserID(ctx),
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
		auth.MustUserID(ctx),
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
