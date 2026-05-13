package handler

import (
	"example/project/backend/models"
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateRoleInput struct {
	Role string `json:"role"`
}

type CreateDoctorInput struct {
	UserID     int64  `json:"user_id"`
	Speciality string `json:"speciality"`
}

type AdminHandler struct {
	analyticsSrv service.AnalyticsService
	inventorySrv service.InventoryService
	usersServ    service.UsersService
	doctorSrv    service.DoctorService // <-- Добавь это
}

func NewAdminHandler(as service.AnalyticsService, is service.InventoryService, us service.UsersService, ds service.DoctorService) *AdminHandler {
	return &AdminHandler{
		analyticsSrv: as,
		inventorySrv: is,
		usersServ:    us,
		doctorSrv:    ds, // Теперь у админа есть доступ к сервису врачей
	}
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

// GetAllUsers godoc
// @Summary      Получение всех пользователей из бд
// @Tags         admin
// @Produce      json
// @Success      200 {array}  models.User
// @Failure      500 {object} map[string]string
// @Router       /api/admin/users [get]
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	users, err := h.usersServ.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список пользователей"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByRole godoc
// @Summary      Получение пользователя по id и роли
// @Tags         admin
// @Produce      json
// @Param        id   path int    true "ID пользователя"
// @Param        role path string true "Роль пользователя"
// @Success      200 {object} models.User
// @Failure      404 {object} map[string]string
// @Router       /api/admin/users/{id}/{role} [get]
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	err = h.usersServ.UpdateUserRole(id, models.UserRole(input.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Роль обновлена"})
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

func (h *AdminHandler) CreateDoctor(c *gin.Context) {
	var input CreateDoctorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	newDoctor := &models.Doctor{
		UserID:     input.UserID,
		Speciality: input.Speciality,
	}

	if err := h.doctorSrv.CreateDoctor(newDoctor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Профиль врача создан"})
}
