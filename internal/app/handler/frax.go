package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /api/frax/cart - иконка корзины
func (h *Handler) APIGetCartBadge(c *gin.Context) {
	draft, err := h.Repository.GetDraftFrax(hardcodedUserID)
	if err != nil {
		c.JSON(http.StatusOK, ds.CartBadgeDTO{
			FraxID: nil,
			Count:  0,
		})
		return
	}

	fullFrax, err := h.Repository.GetFraxWithFactors(draft.ID)
	if err != nil {
		logrus.Error("Error getting frax with factors:", err)
		c.JSON(http.StatusOK, ds.CartBadgeDTO{
			FraxID: nil,
			Count:  0,
		})
		return
	}

	c.JSON(http.StatusOK, ds.CartBadgeDTO{
		FraxID: &fullFrax.ID,
		Count:  len(fullFrax.FactorsLink),
	})
}

// GET /api/frax - список заявок с фильтрацией
func (h *Handler) APIListFrax(c *gin.Context) {
	status := c.Query("status")
	from := c.Query("from") // YYYY-MM-DD
	to := c.Query("to")

	fraxList, err := h.Repository.FraxListFiltered(status, from, to)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fraxList)
}

// GET /api/frax/:id - одна заявка с услугами
func (h *Handler) APIGetFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	frax, err := h.Repository.GetFraxWithFactors(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	var factors []ds.FactorInFraxDTO
	for _, link := range frax.FactorsLink {
		factors = append(factors, ds.FactorInFraxDTO{
			ID:          link.ID,
			FactorID:    link.FactorID,
			Title:       link.Factor.Title,
			Text:        link.Factor.Text,
			Image:       link.Factor.Image,
			Argument:    link.Factor.Argument,
			Description: link.Description,
		})
	}

	fraxDTO := ds.FraxDTO{
		ID:             frax.ID,
		Status:         frax.Status,
		CreationDate:   frax.CreationDate,
		CreatorLogin:   frax.Creator.Username,
		ModeratorLogin: nil,
		FormingDate:    frax.FormingDate,
		ComplitionDate: frax.ComplitionDate,
		Age:            frax.Age,
		Gender:         frax.Gender,
		Weight:         frax.Weight,
		Height:         frax.Height,
		POF:            frax.POF,
		PHF:            frax.PHF,
		Factors:        factors,
	}

	if frax.ModeratorID != nil {
		fraxDTO.ModeratorLogin = &frax.Moderator.Username
	}

	c.JSON(http.StatusOK, fraxDTO)
}

// PUT /api/frax/:id - изменение полей заявки
func (h *Handler) APIUpdateFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.FraxUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateFraxUserFields(uint(id), req); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// PUT /api/frax/:id/form - сформировать заявку
func (h *Handler) APIFormFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.FormFrax(uint(id), hardcodedUserID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// PUT /api/frax/:id/resolve - завершить/отклонить заявку
func (h *Handler) APIResolveFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.FraxResolveRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Для лабораторной работы используем модератора с ID=2
	moderatorID := uint(2)
	if err := h.Repository.ResolveFrax(uint(id), moderatorID, req.Action); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DELETE /api/frax/:id - удаление заявки
func (h *Handler) APIDeleteFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.LogicallyDeleteFrax(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DELETE /api/frax/:id/factors/:factor_id - удаление фактора из заявки
func (h *Handler) APIRemoveFactorFromFrax(c *gin.Context) {
	fraxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	factorID, err := strconv.Atoi(c.Param("factor_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.RemoveFactorFromFrax(uint(fraxID), uint(factorID)); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// PUT /api/frax/:id/factors/:factor_id - изменение м-м связи
func (h *Handler) APIUpdateMM(c *gin.Context) {
	fraxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	factorID, err := strconv.Atoi(c.Param("factor_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.FactorToFraxUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	updateData := ds.FactorToFrax{
		Description: req.Description,
	}

	if err := h.Repository.UpdateMM(uint(fraxID), uint(factorID), updateData); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}
