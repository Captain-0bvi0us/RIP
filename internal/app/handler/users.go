package handler

import (
	"RIP/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/users - регистрация пользователя
func (h *Handler) APIRegister(c *gin.Context) {
	var req ds.UserRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	user := ds.Users{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Moderator: false,
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	userDTO := ds.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Moderator: user.Moderator,
	}

	c.JSON(http.StatusCreated, userDTO)
}

// POST /api/auth/login - аутентификация
func (h *Handler) APILogin(c *gin.Context) {
	var req ds.UserLoginRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByUsername(req.Username)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}

	// Для лабораторной работы возвращаем заглушку токена
	response := ds.LoginResponse{
		Token: "stub_token_for_lab",
		User: ds.UserDTO{
			ID:        user.ID,
			Username:  user.Username,
			Moderator: user.Moderator,
		},
	}

	c.JSON(http.StatusOK, response)
}

// POST /api/auth/logout - деавторизация
func (h *Handler) APILogout(c *gin.Context) {
	// Для лабораторной работы просто возвращаем успех
	c.Status(http.StatusNoContent)
}

// GET /api/users/me - получение данных пользователя
func (h *Handler) APIGetMe(c *gin.Context) {
	// Для лабораторной работы используем захардкоженного пользователя
	user, err := h.Repository.GetUserByID(hardcodedUserID)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	userDTO := ds.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Moderator: user.Moderator,
	}

	c.JSON(http.StatusOK, userDTO)
}

// PUT /api/users/me - обновление данных пользователя
func (h *Handler) APIUpdateMe(c *gin.Context) {
	var req ds.UserUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Для лабораторной работы используем захардкоженного пользователя
	if err := h.Repository.UpdateUser(hardcodedUserID, req); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
