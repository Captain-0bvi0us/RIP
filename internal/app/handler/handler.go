package handler

import (
	"RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const hardcodedUserID = 1

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterAPI(r *gin.RouterGroup) {
	// Домен услуг (факторов)
	r.GET("/factors", h.GetFactors)
	r.GET("/factors/:id", h.GetFactor)
	r.POST("/factors", h.CreateFactor)
	r.PUT("/factors/:id", h.UpdateFactor)
	r.DELETE("/factors/:id", h.DeleteFactor)
	r.POST("/frax/draft/factors/:factor_id", h.AddFactorToDraft)
	r.POST("/factors/:id/image", h.UploadFactorImage)

	// Домен заявок (FRAX)
	r.GET("/frax/cart", h.GetCartBadge)
	r.GET("/frax", h.ListFrax)
	r.GET("/frax/:id", h.GetFrax)
	r.PUT("/frax/:id", h.UpdateFrax)
	r.PUT("/frax/:id/form", h.FormFrax)
	r.PUT("/frax/:id/resolve", h.ResolveFrax)
	r.DELETE("/frax/:id", h.DeleteFrax)

	// Домен м-м
	r.DELETE("/frax/:id/factors/:factor_id", h.RemoveFactorFromFrax)
	r.PUT("/frax/:id/factors/:factor_id", h.UpdateMM)

	// Домен пользователь
	r.POST("/users", h.Register)
	r.GET("/users/:id", h.GetUserData)
	r.PUT("/users/:id", h.UpdateUserData)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
