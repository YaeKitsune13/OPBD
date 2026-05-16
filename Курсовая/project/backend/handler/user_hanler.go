package handler

import (
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSrv service.UsersService
}

func NewUserHandler(userSrv service.UsersService) *UserHandler {
	return &UserHandler{userSrv: userSrv}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var input struct {
		LastName   string `json:"lastName"`
		FirstName  string `json:"firstName"`
		MiddleName string `json:"middleName"`
		Phone      string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	err := h.userSrv.UpdateUser(id, input.LastName, input.FirstName, input.MiddleName, input.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Профиль обновлён"})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var input struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	err := h.userSrv.ChangePassword(id, input.CurrentPassword, input.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль изменён"})
}
