package handler

import (
	"RIP/internal/app/ds"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /api/frax/factorscart - иконка корзины

// GetCartBadge godoc
// @Summary      Получить информацию для иконки корзины (авторизованный пользователь)
// @Description  Возвращает ID черновика текущего пользователя и количество факторов в нем.
// @Tags         frax
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} ds.CartBadgeDTO
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /frax/factorscart [get]
func (h *Handler) GetCartBadge(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}

	draft, err := h.Repository.GetDraftFrax(userID)
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

// ListFrax godoc
// @Summary      Получить список заявок (авторизованный пользователь)
// @Description  Возвращает отфильтрованный список всех сформированных заявок (кроме черновиков и удаленных).
// @Tags         frax
// @Produce      json
// @Security     ApiKeyAuth
// @Param        status query int false "Фильтр по статусу заявки"
// @Param        from query string false "Фильтр по дате 'от' (формат YYYY-MM-DD)"
// @Param        to query string false "Фильтр по дате 'до' (формат YYYY-MM-DD)"
// @Success      200 {array} ds.FraxDTO
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /frax [get]
func (h *Handler) ListFrax(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}
	isModerator := isUserModerator(c)

	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")

	fraxList, err := h.Repository.FraxListFiltered(userID, isModerator, status, from, to)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fraxList)
}

// GET /api/frax/:id - одна заявка с услугами

// GetFrax godoc
// @Summary      Получить одну заявку по ID (авторизованный пользователь)
// @Description  Возвращает полную информацию о заявке, включая привязанные факторы.
// @Tags         frax
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Success      200 {object} ds.FraxDTO
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      404 {object} map[string]string "Заявка не найдена"
// @Router       /frax/{id} [get]
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

// UpdateFrax godoc
// @Summary      Обновить данные заявки (авторизованный пользователь)
// @Description  Позволяет пользователю обновить поля своей заявки (возраст, пол, вес, рост).
// @Tags         frax
// @Accept       json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        updateData body ds.FraxUpdateRequest true "Данные для обновления"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /frax/{id} [put]
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

// DELETE /api/frax/:id - удаление заявки

// DeleteFrax godoc
// @Summary      Удалить заявку (авторизованный пользователь)
// @Description  Логически удаляет заявку, переводя ее в статус "удалена".
// @Tags         frax
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /frax/{id} [delete]
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

// RemoveFactorFromFrax godoc
// @Summary      Удалить фактор из заявки (авторизованный пользователь)
// @Description  Удаляет связь между заявкой и фактором.
// @Tags         m-m
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        factor_id path int true "ID фактора"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /frax/{id}/factors/{factor_id} [delete]
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

// UpdateMM godoc
// @Summary      Обновить описание фактора в заявке (авторизованный пользователь)
// @Description  Изменяет дополнительное описание для конкретного фактора в рамках одной заявки.
// @Tags         m-m
// @Accept       json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        factor_id path int true "ID фактора"
// @Param        updateData body ds.FactorToFraxUpdateRequest true "Новое описание"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /frax/{id}/factors/{factor_id} [put]
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

// PUT /api/internal/frax/result
func (h *Handler) SetFraxResult(c *gin.Context) {
	token := c.GetHeader("Authorization")
	expectedToken := "secret12"

	if token != expectedToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}

	var res ds.AsyncCalcResponse
	if err := c.BindJSON(&res); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateFraxResults(res.ID, res.POF, res.PHF); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "results updated"})
}

// PUT /api/frax/:id/form - сформировать заявку

// FormFrax godoc
// @Summary      Сформировать заявку (авторизованный пользователь)
// @Description  Переводит заявку из статуса "черновик" в "сформирована".
// @Tags         frax
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки (черновика)"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Не все поля заполнены"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /frax/{id}/form [put]

func (h *Handler) FormFrax(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}

	if err := h.Repository.FormFrax(uint(id), userID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка сформирована и отправлена модератору",
	})
}

// PUT /api/frax/:id/resolve - завершить/отклонить заявку

// ResolveFrax godoc
// @Summary      Завершить или отклонить заявку (только модератор)
// @Description  Модератор завершает (с расчетом) или отклоняет заявку.
// @Tags         frax
// @Accept       json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        action body ds.FraxResolveRequest true "Действие: 'complete' или 'reject'"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /frax/{id}/resolve [put]

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

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}
	moderatorID := uint(userID)

	if err := h.Repository.ResolveFrax(uint(id), moderatorID, req.Action); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if req.Action == "complete" {
		fraxFull, err := h.Repository.GetFraxWithFactors(uint(id))
		if err == nil {
			factorSum := 0.0
			for _, link := range fraxFull.FactorsLink {
				if link.Factor.Argument != nil {
					factorSum += *link.Factor.Argument
				}
			}

			reqData := ds.AsyncCalcRequest{
				ID:        fraxFull.ID,
				Age:       *fraxFull.Age,
				Gender:    *fraxFull.Gender,
				Weight:    *fraxFull.Weight,
				Height:    *fraxFull.Height,
				FactorSum: factorSum,
			}

			go sendAsyncCalculation("http://localhost:8000/api/calculate_probability/", reqData)
		} else {
			logrus.Errorf("Failed to fetch frax data for async calc: %v", err)
		}
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка обработана модератором",
	})
}

func sendAsyncCalculation(url string, data ds.AsyncCalcRequest) {
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorf("Failed to send async calc request: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Async service returned non-200 status: %d", resp.StatusCode)
	}
}
