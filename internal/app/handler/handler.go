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

// Регистрация только API роутов
func (h *Handler) RegisterAPI(r *gin.RouterGroup) {
	// Домен услуг (факторов)
	r.GET("/factors", h.GetFactors)
	r.GET("/factors/:id", h.GetFactor)
	r.POST("/factors", h.CreateFactor)
	r.PUT("/factors/:id", h.UpdateFactor)
	r.DELETE("/factors/:id", h.DeleteFactor)
	r.POST("/frax/draft/factors/:factor_id", h.APIAddFactorToDraft)
	r.POST("/factors/:id/image", h.UploadFactorImage)

	// Домен заявок (FRAX)
	r.GET("/frax/cart", h.APIGetCartBadge)
	r.GET("/frax", h.APIListFrax)
	r.GET("/frax/:id", h.APIGetFrax)
	r.PUT("/frax/:id", h.APIUpdateFrax)
	r.PUT("/frax/:id/form", h.APIFormFrax)
	r.PUT("/frax/:id/resolve", h.APIResolveFrax)
	r.DELETE("/frax/:id", h.APIDeleteFrax)

	// Домен м-м
	r.DELETE("/frax/:id/factors/:factor_id", h.APIRemoveFactorFromFrax)
	r.PUT("/frax/:id/factors/:factor_id", h.APIUpdateMM)

	// Домен пользователь
	r.POST("/users", h.APIRegister)
	r.POST("/auth/login", h.APILogin)
	r.POST("/auth/logout", h.APILogout)
	r.GET("/users/me", h.APIGetMe)
	r.PUT("/users/me", h.APIUpdateMe)
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
