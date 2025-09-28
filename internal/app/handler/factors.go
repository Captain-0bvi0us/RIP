package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/factors - список факторов с фильтрацией
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

	c.Status(http.StatusNoContent)
}

// POST /api/frax/draft/factors/:factor_id - добавление фактора в черновик
func (h *Handler) AddFactorToDraft(c *gin.Context) {
	factorID, err := strconv.Atoi(c.Param("factor_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.AddFactorToDraft(hardcodedUserID, uint(factorID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

// POST /api/factors/:id/image - загрузка изображения фактора
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
