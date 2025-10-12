package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/factors - список факторов с фильтрацией

// GetFactors godoc
// @Summary      Получить список факторов
// @Description  Возвращает постраничный список факторов риска. Доступен для всех авторизованных пользователей.
// @Tags         factors
// @Produce      json
// @Security     ApiKeyAuth
// @Param        title query string false "Фильтр по названию фактора (поиск по подстроке)"
// @Success      200 {object} ds.PaginatedResponse
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /factors [get]
func (h *Handler) GetFactors(c *gin.Context) {
	title := c.Query("title")

	factors, total, err := h.Repository.FactorsList(title)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var factorDTOs []ds.FactorDTO
	for _, f := range factors {
		factorDTOs = append(factorDTOs, ds.FactorDTO{
			ID:       f.ID,
			Title:    f.Title,
			Text:     f.Text,
			Image:    f.Image,
			Argument: f.Argument,
			Status:   f.Status,
		})
	}

	c.JSON(http.StatusOK, ds.PaginatedResponse{
		Items: factorDTOs,
		Total: total,
	})
}

// GET /api/factors/:id - один фактор

// GetFactor godoc
// @Summary      Получить один фактор по ID
// @Description  Возвращает детальную информацию о факторе риска.
// @Tags         factors
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "ID фактора"
// @Success      200 {object} ds.FactorDTO
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      404 {object} map[string]string "Фактор не найден"
// @Router       /factors/{id} [get]
func (h *Handler) GetFactor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	factor, err := h.Repository.GetFactorByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	factorDTO := ds.FactorDTO{
		ID:       factor.ID,
		Title:    factor.Title,
		Text:     factor.Text,
		Image:    factor.Image,
		Argument: factor.Argument,
		Status:   factor.Status,
	}

	c.JSON(http.StatusOK, factorDTO)
}

// POST /api/factors - создание фактора

// CreateFactor godoc
// @Summary      Создать новый фактор (только модератор)
// @Description  Создает новую запись о факторе риска. Доступно только для модераторов.
// @Tags         factors
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        factorData body ds.FactorCreateRequest true "Данные нового фактора"
// @Success      201 {object} ds.FactorDTO
// @Failure      400 {object} map[string]string "Ошибка валидации"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен (не модератор)"
// @Router       /factors [post]
func (h *Handler) CreateFactor(c *gin.Context) {
	var req ds.FactorCreateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	statusValue := false

	factor := ds.Factors{
		Title:    req.Title,
		Text:     req.Text,
		Argument: req.Argument,
		Status:   &statusValue,
	}

	if err := h.Repository.CreateFactor(&factor); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	factorDTO := ds.FactorDTO{
		ID:       factor.ID,
		Title:    factor.Title,
		Text:     factor.Text,
		Image:    factor.Image,
		Argument: factor.Argument,
		Status:   factor.Status,
	}

	c.JSON(http.StatusCreated, factorDTO)
}

// PUT /api/factors/:id - обновление фактора

// UpdateFactor godoc
// @Summary      Обновить фактор (только модератор)
// @Description  Обновляет информацию о существующем факторе риска.
// @Tags         factors
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "ID фактора"
// @Param        updateData body ds.FactorUpdateRequest true "Данные для обновления"
// @Success      200 {object} ds.FactorDTO
// @Failure      400 {object} map[string]string "Ошибка валидации"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /factors/{id} [put]
func (h *Handler) UpdateFactor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.FactorUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	factor, err := h.Repository.UpdateFactor(uint(id), req)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	factorDTO := ds.FactorDTO{
		ID:       factor.ID,
		Title:    factor.Title,
		Text:     factor.Text,
		Image:    factor.Image,
		Argument: factor.Argument,
		Status:   factor.Status,
	}

	c.JSON(http.StatusOK, factorDTO)
}

// DELETE /api/factors/:id - удаление фактора

// DeleteFactor godoc
// @Summary      Удалить фактор (только модератор)
// @Description  Удаляет фактор риска из системы.
// @Tags         factors
// @Security     ApiKeyAuth
// @Param        id path int true "ID фактора для удаления"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /factors/{id} [delete]
func (h *Handler) DeleteFactor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.DeleteFactor(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Фактор удален",
	})
}

// POST /api/frax/draft/factors/:factor_id - добавление фактора в черновик

// AddFactorToDraft godoc
// @Summary      Добавить фактор в черновик заявки
// @Description  Находит или создает черновик заявки для текущего пользователя и добавляет в него фактор.
// @Tags         frax
// @Security     ApiKeyAuth
// @Param        factor_id path int true "ID фактора для добавления"
// @Success      201 {object} map[string]string "Сообщение об успехе"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /frax/draft/factors/{factor_id} [post]
func (h *Handler) AddFactorToDraft(c *gin.Context) {
	factorID, err := strconv.Atoi(c.Param("factor_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}

	if err := h.Repository.AddFactorToDraft(userID, uint(factorID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Черновик создан. Фактор добавлен в черновик.",
	})
}

// POST /api/factors/:id/image - загрузка изображения фактора

// UploadFactorImage godoc
// @Summary      Загрузить изображение для фактора (только модератор)
// @Description  Загружает и привязывает изображение к фактору риска.
// @Tags         factors
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "ID фактора"
// @Param        file formData file true "Файл изображения"
// @Success      200 {object} map[string]string "URL загруженного изображения"
// @Failure      400 {object} map[string]string "Файл не предоставлен"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /factors/{id}/image [post]
func (h *Handler) UploadFactorImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	imageURL, err := h.Repository.UploadFactorImage(uint(id), file)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": imageURL})
}
