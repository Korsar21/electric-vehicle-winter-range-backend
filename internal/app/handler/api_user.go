package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/serializer"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *APIHandler) RegisterUser(ctx *gin.Context) {
	var request serializer.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "fail",
				"message": "invalid request body",
			},
		)
		return
	}

	request.Name = strings.TrimSpace(request.Name)

	if request.Name == "" || request.Password == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "fail",
				"message": "name and password are required",
			},
		)
		return
	}

	_, err := h.Repository.GetUserByName(request.Name)

	if err == nil {
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "fail",
				"message": "user already exists",
			},
		)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "fail",
				"message": "failed to check user",
			},
		)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "fail",
				"message": "failed to hash password",
			},
		)
		return
	}

	user := model.User{
		Name:     request.Name,
		Password: string(hash),
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "fail",
				"message": "failed to create user",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		serializer.UserResponse{
			ID:   user.ID,
			Name: user.Name,
		},
	)
}

func (h *APIHandler) AuthenticateUser(ctx *gin.Context) {
	var request serializer.UserAuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "fail",
				"message": "invalid request body",
			},
		)
		return
	}

	request.Name = strings.TrimSpace(request.Name)

	user, err := h.Repository.GetUserByName(request.Name)

	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status":  "fail",
				"message": "invalid credentials",
			},
		)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status":  "fail",
				"message": "invalid credentials",
			},
		)
		return
	}

	token, expiresAt, err := h.Auth.GenerateToken(
		user.ID,
		user.Name,
	)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "fail",
				"message": "failed to create token",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializer.UserAuthResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
			UserID:      user.ID,
			Name:        user.Name,
			Role:        "creator",
		},
	)
}

func (h *APIHandler) LogoutUser(ctx *gin.Context) {
	if err := h.Auth.Logout(ctx); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status":  "fail",
				"message": "failed to revoke token",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"status": "success",
			"data": gin.H{
				"logged_out": true,
			},
		},
	)
}
