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

	ctx.HTML(http.StatusOK, "factors.html", gin.H{
		"factors":       factors,
		"factorsSearch": search,
	})
}

func (h *Handler) GetFactorById(ctx *gin.Context) {
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

// func (h *Handler) GetFactor(ctx *gin.Context) {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	factor, err := h.Repository.GetFactor(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	ctx.HTML(http.StatusOK, "oneFactor.html", gin.H{
// 		"factor": factor,
// 	})
// }

// // frax

// func (h *Handler) GetFraxPage(ctx *gin.Context) {
// 	idStr := ctx.Param("id")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	fraxPage, err := h.Repository.GetFraxPage(id)
// 	factorsInFraxPage := fraxPage.Factors
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	ctx.HTML(http.StatusOK, "fraxPage.html", gin.H{
// 		"factorsInFraxPage": factorsInFraxPage,
// 		"Age":               fraxPage.Age,
// 		"Gender":            fraxPage.Gender,
// 		"Weight":            fraxPage.Weight,
// 		"Height":            fraxPage.Height,
// 		"FirstResult":       fraxPage.FirstResult,
// 		"SecondResult":      fraxPage.SecondResult,
// 	})
// }
