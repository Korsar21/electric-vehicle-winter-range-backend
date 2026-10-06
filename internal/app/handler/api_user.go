package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/serializer"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /api/users/register
//
// егистрация нового пользователя.
func (h *APIHandler) RegisterUser(ctx *gin.Context) {
	var request serializer.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid request body",
		})
		return
	}

	request.Name = strings.TrimSpace(request.Name)

	if request.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "name is required",
		})
		return
	}

	if request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "password is required",
		})
		return
	}

	_, err := h.Repository.GetUserByName(request.Name)

	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"status":  "fail",
			"message": "user already exists",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to check user",
		})
		return
	}

	user := model.User{
		Name:     request.Name,
		Password: request.Password,
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to create user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, serializer.UserResponse{
		ID:   user.ID,
		Name: user.Name,
	})
}

// POST /api/users/auth
//
// аглушка аутентификации для 3.
// олноценная авторизация реализуется в 4.
func (h *APIHandler) AuthenticateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "authentication stub",
	})
}

// POST /api/users/logout
//
// аглушка деавторизации для 3.
// олноценная авторизация реализуется в 4.
func (h *APIHandler) LogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logout stub",
	})
}
