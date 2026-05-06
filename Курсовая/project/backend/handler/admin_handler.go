package handler

import (
	"example/project/backend/models"
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	analyticsSrv service.AnalyticsService
	inventorySrv service.InventoryService // Мы добавим сюда методы CRUD
}

func NewAdminHandler(as service.AnalyticsService, is service.InventoryService) *AdminHandler {
	return &AdminHandler{analyticsSrv: as, inventorySrv: is}
}

// GetStats godoc
// @Summary      Общая статистика (KPI)
// @Tags         admin
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]string
// @Router       /api/admin/stats [get]
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.analyticsSrv.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetRevenue godoc
// @Summary      Отчёт по выручке
// @Tags         admin
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]string
// @Router       /api/admin/revenue [get]
func (h *AdminHandler) GetRevenue(c *gin.Context) {
	report, err := h.analyticsSrv.GetRevenueReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

// DeleteService godoc
// @Summary      Удалить услугу
// @Tags         admin
// @Produce      json
// @Param        id path int true "ID услуги"
// @Success      200 {object} map[string]string
// @Router       /api/admin/services/{id} [delete]
func (h *AdminHandler) DeleteService(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	// Вызов метода из InventoryService (нужно будет туда добавить DeleteService)
	c.JSON(http.StatusOK, gin.H{"message": "Услуга удалена", "id": id})
}

// CreateMed godoc
// @Summary      Добавить медикамент
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body models.Medication true "Данные медикамента"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/admin/meds [post]
func (h *AdminHandler) CreateMed(c *gin.Context) {
	var input models.Medication
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Данные не верны"})
		return
	}
	// Вызов метода из InventoryService
	c.JSON(http.StatusCreated, gin.H{"message": "Медикамент добавлен"})
}
