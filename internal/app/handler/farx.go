package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const hardcodedUserID = 1 // Хардкодим ID пользователя, как указано в задаче

// AddFactorToFrax добавляет фактор в заявку-черновик
func (h *Handler) AddFactorToFrax(c *gin.Context) {
	factorID, err := strconv.Atoi(c.Param("factor_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// 1. Найти или создать черновик для пользователя
	frax, err := h.Repository.GetOrCreateDraftFrax(hardcodedUserID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// 2. Добавить фактор в этот черновик
	err = h.Repository.AddFactorToFrax(frax.ID, uint(factorID))
	if err != nil {
		// Можно проигнорировать ошибку "уже существует", если это приемлемо
	}

	// 3. Перенаправить пользователя обратно на страницу с факторами
	c.Redirect(http.StatusFound, "/FRAX")
}

// GetFraxPage отображает страницу составления заявки
func (h *Handler) GetFraxPage(c *gin.Context) {
	fraxID, err := strconv.Atoi(c.Param("frax_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Получаем заявку со всеми факторами
	frax, err := h.Repository.GetFraxWithFactors(uint(fraxID))
	if err != nil {
		// Если заявка не найдена или удалена, показываем ошибку
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Невозможно зайти на страницу при отсутствии факторов
	if len(frax.FactorsLink) == 0 {
		h.errorHandler(c, http.StatusForbidden, errors.New("cannot access an empty frax page, add factors first"))
		return
	}

	c.HTML(http.StatusOK, "fraxPage.html", frax)
}

// DeleteFrax логически удаляет всю заявку
func (h *Handler) DeleteFrax(c *gin.Context) {
	fraxID, _ := strconv.Atoi(c.Param("frax_id"))

	if err := h.Repository.LogicallyDeleteFrax(uint(fraxID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Перенаправляем на главную страницу с факторами
	c.Redirect(http.StatusFound, "/FRAX")
}
