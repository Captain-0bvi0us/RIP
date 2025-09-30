package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /api/frax/cart - иконка корзины
func (h *Handler) GetCartBadge(c *gin.Context) {
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
func (h *Handler) ListFrax(c *gin.Context) {
	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")

	fraxList, err := h.Repository.FraxListFiltered(status, from, to)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fraxList)
}

// GET /api/frax/:id - одна заявка с услугами
func (h *Handler) GetFrax(c *gin.Context) {
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
		CreatorID:      frax.Creator.ID,
		ModeratorID:    nil,
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
		fraxDTO.ModeratorID = &frax.Moderator.ID
	}

	c.JSON(http.StatusOK, fraxDTO)
}

// PUT /api/frax/:id - изменение полей заявки
func (h *Handler) UpdateFrax(c *gin.Context) {
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

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Данные заявки обновлены",
	})
}

// PUT /api/frax/:id/form - сформировать заявку
func (h *Handler) FormFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.FormFrax(uint(id), hardcodedUserID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка сформирована",
	})
}

// PUT /api/frax/:id/resolve - завершить/отклонить заявку
func (h *Handler) ResolveFrax(c *gin.Context) {
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

	moderatorID := uint(hardcodedUserID)
	if err := h.Repository.ResolveFrax(uint(id), moderatorID, req.Action); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка обработана модератором",
	})
}

// DELETE /api/frax/:id - удаление заявки
func (h *Handler) DeleteFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.LogicallyDeleteFrax(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка удалена",
	})
}

// DELETE /api/frax/:id/factors/:factor_id - удаление фактора из заявки
func (h *Handler) RemoveFactorFromFrax(c *gin.Context) {
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

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Фактор удален из заявки",
	})
}

// PUT /api/frax/:id/factors/:factor_id - изменение м-м связи
func (h *Handler) UpdateMM(c *gin.Context) {
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

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Дополнительная информация к фаткору обновлена",
	})
}
