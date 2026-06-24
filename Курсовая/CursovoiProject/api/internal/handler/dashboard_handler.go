package handler

import (
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc}
}

// GetDashboard godoc
// @Summary      Данные дашборда
// @Description  Получает сводку статистики, ближайшую запись и список питомцев
// @Tags         dashboard
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "ID пользователя"
// @Success      200     {object}  dto.DashboardResponse
// @Router       /dashboard/{userId} [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	data, err := h.svc.GetDashboardData(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load dashboard"})
		return
	}
	c.JSON(http.StatusOK, data)
}
