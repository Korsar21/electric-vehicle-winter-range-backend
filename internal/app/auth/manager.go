package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const contextUserIDKey = "auth_user_id"

type Claims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
	ttl    time.Duration
	redis  *redis.Client
}

func NewManager() (*Manager, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "lr4-development-secret-change-me"
	}

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{Addr: addr})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Manager{
		secret: []byte(secret),
		ttl:    time.Hour,
		redis:  client,
	}, nil
}

func (m *Manager) GenerateToken(userID uint, name string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)

	return signed, expiresAt, err
}

func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	return parts[1], true
}

func (m *Manager) isBlacklisted(ctx context.Context, jti string) (bool, error) {
	count, err := m.redis.Exists(ctx, "jwt:blacklist:"+jti).Result()
	return count > 0, err
}

func (m *Manager) blacklist(ctx context.Context, claims *Claims) error {
	if claims == nil || claims.ID == "" || claims.ExpiresAt == nil {
		return errors.New("token claims are incomplete")
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		ttl = time.Minute
	}

	value := fmt.Sprintf("user_id=%d;name=%s", claims.UserID, claims.Name)

	return m.redis.Set(
		ctx,
		"jwt:blacklist:"+claims.ID,
		value,
		ttl,
	).Err()
}

func UserIDFromContext(ctx *gin.Context) (uint, bool) {
	value, ok := ctx.Get(contextUserIDKey)
	if !ok {
		return 0, false
	}

	userID, ok := value.(uint)

	return userID, ok && userID > 0
}

func MustUserID(ctx *gin.Context) uint {
	userID, _ := UserIDFromContext(ctx)
	return userID
}

func (m *Manager) Required() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, ok := bearerToken(ctx.GetHeader("Authorization"))

		if !ok {
			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"status":  "fail",
					"message": "authorization token is required",
				},
			)
			return
		}

		claims, err := m.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"status":  "fail",
					"message": "invalid or expired token",
				},
			)
			return
		}

		blocked, err := m.isBlacklisted(
			ctx.Request.Context(),
			claims.ID,
		)

		if err != nil {
			ctx.AbortWithStatusJSON(
				500,
				gin.H{
					"status":  "fail",
					"message": "failed to check token blacklist",
				},
			)
			return
		}

		if blocked {
			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"status":  "fail",
					"message": "token has been revoked",
				},
			)
			return
		}

		ctx.Set(contextUserIDKey, claims.UserID)

		ctx.Next()
	}
}

func (m *Manager) Optional() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, ok := bearerToken(ctx.GetHeader("Authorization"))

		if !ok {
			ctx.Next()
			return
		}

		claims, err := m.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		blocked, err := m.isBlacklisted(
			ctx.Request.Context(),
			claims.ID,
		)

		if err == nil && !blocked {
			ctx.Set(contextUserIDKey, claims.UserID)
		}

		ctx.Next()
	}
}

func (m *Manager) Logout(ctx *gin.Context) error {
	tokenString, ok := bearerToken(ctx.GetHeader("Authorization"))

	if !ok {
		return errors.New("authorization token is required")
	}

	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return err
	}

	return m.blacklist(
		ctx.Request.Context(),
		claims,
	)
}
