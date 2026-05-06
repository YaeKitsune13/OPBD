package handler

import (
	"net/http"
	"strconv"

	"example/project/backend/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashSrv service.DashboardService
}

func NewDashboardHandler(dashSrv service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashSrv: dashSrv}
}

// GetData godoc
// @Summary      Данные дашборда владельца
// @Tags         dashboard
// @Produce      json
// @Param        ownerId path int true "ID владельца"
// @Success      200 {object} dto.DashboardDTO
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/dashboard/{ownerId} [get]
func (h *DashboardHandler) GetData(c *gin.Context) {
	ownerIdStr := c.Param("ownerId")
	ownerId, err := strconv.ParseInt(ownerIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ownerId"})
		return
	}

	data, err := h.dashSrv.GetDashboardData(ownerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
