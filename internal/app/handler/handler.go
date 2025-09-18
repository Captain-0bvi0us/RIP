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
	FraxID := 1

	frax, err := h.Repository.GetFrax(FraxID)
	if err != nil {
		logrus.Error(err)
	}
	FactorsCount := len(frax.Factors)

	searchFactor := ctx.Query("searchingFactors")
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
		"fraxID":       FraxID,
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

func (h *Handler) GetFrax(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	frax, err := h.Repository.GetFrax(id)
	factorsInFrax := frax.Factors
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "frax.html", gin.H{
		"factorsInFrax": factorsInFrax,
		"Age":           frax.Age,
		"Gender":        frax.Gender,
		"Weight":        frax.Weight,
		"Height":        frax.Height,
		"POF":           frax.POF,
		"PHF":           frax.PHF,
	})
}
