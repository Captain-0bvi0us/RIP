package handler

import (
	"RIP/internal/app/config"
	"RIP/internal/app/redis"
	"RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	Redis      *redis.Client
	JWTConfig  *config.JWTConfig
}

func NewHandler(r *repository.Repository, redis *redis.Client, jwtConfig *config.JWTConfig) *Handler {
	return &Handler{
		Repository: r,
		Redis:      redis,
		JWTConfig:  jwtConfig,
	}
}

func (h *Handler) RegisterAPI(r *gin.RouterGroup) {

	// Доступны всем
	r.POST("/users", h.Register)
	r.POST("/auth/login", h.Login)
	r.GET("/factors", h.GetFactors)
	r.GET("/factors/:id", h.GetFactor)

	// Эндпоинты, доступные только авторизованным пользователям
	auth := r.Group("/")
	auth.Use(h.AuthMiddleware)
	{
		// Пользователи
		auth.POST("/auth/logout", h.Logout)
		auth.GET("/users/:id", h.GetUserData)
		auth.PUT("/users/:id", h.UpdateUserData)

		// Заявки
		auth.POST("/frax/draft/factors/:factor_id", h.AddFactorToDraft)
		auth.GET("/frax/factorscart", h.GetCartBadge)
		auth.GET("/frax", h.ListFrax)
		auth.GET("/frax/:id", h.GetFrax)
		auth.PUT("/frax/:id", h.UpdateFrax)
		auth.PUT("/frax/:id/form", h.FormFrax)
		auth.DELETE("/frax/:id", h.DeleteFrax)
		auth.DELETE("/frax/:id/factors/:factor_id", h.RemoveFactorFromFrax)
		auth.PUT("/frax/:id/factors/:factor_id", h.UpdateMM)
	}

	// Эндпоинты, доступные только модераторам
	moderator := r.Group("/")
	moderator.Use(h.AuthMiddleware, h.ModeratorMiddleware)
	{
		// Управление факторами (создание, изменение, удаление)
		moderator.POST("/factors", h.CreateFactor)
		moderator.PUT("/factors/:id", h.UpdateFactor)
		moderator.DELETE("/factors/:id", h.DeleteFactor)
		moderator.POST("/factors/:id/image", h.UploadFactorImage)

		// Управление заявками (завершение/отклонение)
		moderator.PUT("/frax/:id/resolve", h.ResolveFrax)
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
