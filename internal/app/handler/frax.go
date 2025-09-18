package handler

import (
	"RIP/internal/app/ds"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const hardcodedUserID = 1

func (h *Handler) AddFactorToFrax(c *gin.Context) {
	factorID, err := strconv.Atoi(c.Param("factor_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	frax, err := h.Repository.GetDraftFrax(hardcodedUserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newFrax := ds.FraxSearching{
			CreatorID: hardcodedUserID,
			Status:    ds.StatusDraft,
		}
		if createErr := h.Repository.CreateFrax(&newFrax); createErr != nil {
			h.errorHandler(c, http.StatusInternalServerError, createErr)
			return
		}
		frax = &newFrax
	} else if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if err = h.Repository.AddFactorToFrax(frax.ID, uint(factorID)); err != nil {
	}

	c.Redirect(http.StatusFound, "/FRAX")
}

func (h *Handler) GetFrax(c *gin.Context) {
	fraxID, err := strconv.Atoi(c.Param("frax_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	frax, err := h.Repository.GetFraxWithFactors(uint(fraxID))
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	if len(frax.FactorsLink) == 0 {
		h.errorHandler(c, http.StatusForbidden, errors.New("cannot access an empty frax page, add factors first"))
		return
	}

	c.HTML(http.StatusOK, "frax.html", frax)
}

func (h *Handler) DeleteFrax(c *gin.Context) {
	fraxID, _ := strconv.Atoi(c.Param("frax_id"))

	if err := h.Repository.LogicallyDeleteFrax(uint(fraxID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/FRAX")
}
