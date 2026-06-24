package handler

import (
	"api/internal/dto"
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc}
}

// GetProfile godoc
// @Summary      Профиль пользователя
// @Tags         user
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  dto.UserProfileResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	res, err := h.svc.GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateProfile godoc
// @Summary      Обновить профиль
// @Tags         user
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "ID пользователя"
// @Param        request  body      dto.UpdateProfileRequest  true  "Новые данные"
// @Success      200      {object}  map[string]bool
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.svc.UpdateProfile(uint(id), req)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ChangePassword godoc
// @Summary      Сменить пароль
// @Tags         user
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "ID пользователя"
// @Param        request  body      dto.ChangePasswordRequest  true  "Пароли"
// @Success      200      {object}  map[string]bool
// @Router       /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ChangePassword(uint(id), req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
