package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAllFactors(ctx *gin.Context) {
	var factors []ds.Factors
	var err error

	search := ctx.Query("search")
	if search == "" {
		factors, err = h.Repository.GetAllFactors()
	} else {
		factors, err = h.Repository.SearchFactorsByName(search)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	draftFrax, err := h.Repository.GetDraftFrax(hardcodedUserID)
	var fraxID uint = 0
	var factorsCount int = 0

	if err == nil && draftFrax != nil {
		fullFrax, err := h.Repository.GetFraxWithFactors(draftFrax.ID)
		if err == nil {
			fraxID = fullFrax.ID
			factorsCount = len(fullFrax.FactorsLink)
		}
	}

	ctx.HTML(http.StatusOK, "factors.html", gin.H{
		"factors":       factors,
		"factorsSearch": search,
		"fraxID":        fraxID,
		"factorsCount":  factorsCount,
	})
}

func (h *Handler) GetFactorByID(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	factor, err := h.Repository.GetFactorByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "oneFactor.html", factor)
}
