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
	var factors []repository.Factor
	var err error
	FraxPageID := 1

	fraxPage, err := h.Repository.GetFraxPage(FraxPageID)
	if err != nil {
		logrus.Error(err)
	}
	FactorsCount := len(fraxPage.Factors)

	searchFactor := ctx.Query("query")
	if searchFactor == "" {
		factors, err = h.Repository.GetFactors()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		factors, err = h.Repository.GetFactorsByTitle(searchFactor)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "factors.html", gin.H{
		"factors":      factors,
		"searchFactor": searchFactor,
		"factorsCount": FactorsCount,
		"fraxPageID":   FraxPageID,
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

	ctx.HTML(http.StatusOK, "oneFactor.html", gin.H{
		"factor": factor,
	})
}

// frax

func (h *Handler) GetFraxPage(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	fraxPage, err := h.Repository.GetFraxPage(id)
	factorsInFraxPage := fraxPage.Factors
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "fraxPage.html", gin.H{
		"factorsInFraxPage": factorsInFraxPage,
		"Age":               fraxPage.Age,
		"Gender":            fraxPage.Gender,
		"Weight":            fraxPage.Weight,
		"Height":            fraxPage.Height,
		"POF":               fraxPage.POF,
		"PHF":               fraxPage.PHF,
	})
}
