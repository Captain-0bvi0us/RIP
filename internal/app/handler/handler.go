package handler

import (
	"RIP/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// factors

func (h *Handler) GetFactors(ctx *gin.Context) {
	var factors []repository.Factors
	var err error

	// search

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		factors, err = h.Repository.GetFactors()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		factors, err = h.Repository.GetFactorsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// count

	orders, err := h.Repository.GetOrders()
	if err != nil {
		logrus.Error(err)
	}

	allFactors, err := h.Repository.GetFactors()
	if err != nil {
		logrus.Error(err)
	}

	count := 0
	for _, order := range orders {
		for _, factor := range order.Factors {
			for _, f := range allFactors {
				if factor.ID == f.ID {
					count++
				}
			}
		}
	}

	ctx.HTML(http.StatusOK, "factors.html", gin.H{
		"factors": factors,
		"query":   searchQuery,
		"count":   count,
		"orderID": orders[0].ID_order,
	})
}

func (h *Handler) GetFactor(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	factor, err := h.Repository.GetFactor(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "one_factor.html", gin.H{
		"factor": factor,
	})
}

// orders

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
	})
}
