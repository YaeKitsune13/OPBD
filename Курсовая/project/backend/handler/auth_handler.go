package handler

import (
	"example/project/backend/dto"
	"example/project/backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	srv service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{srv: s}
}

// Login godoc
// @Summary      Вход в систему
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginRequest true "Данные для входа"
// @Success      200 {object} dto.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	res, err := h.srv.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Register godoc
// @Summary      Регистрация нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RegisterRequest true "Данные для регистрации"
// @Success      201 {object} dto.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Заполните все поля"})
		return
	}

	res, err := h.srv.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Logout godoc
// @Summary      Выход из системы
// @Tags         auth
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Успешный выход"})
}
