package handler

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateRoleInput struct {
	Role string `json:"role"`
}

type CreateDoctorInput struct {
	UserID     int64  `json:"user_id"`
	Speciality string `json:"speciality"`
}

type RegisterDoctorInput struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone"`
	Password   string `json:"password" binding:"required"`
	Speciality string `json:"speciality" binding:"required"`
}

type AdminHandler struct {
	analyticsSrv service.AnalyticsService
	inventorySrv service.InventoryService
	usersServ    service.UsersService
	doctorSrv    service.DoctorService
}

func NewAdminHandler(as service.AnalyticsService, is service.InventoryService, us service.UsersService, ds service.DoctorService) *AdminHandler {
	return &AdminHandler{
		analyticsSrv: as,
		inventorySrv: is,
		usersServ:    us,
		doctorSrv:    ds,
	}
}

// DeleteUser godoc
// @Summary      Удалить пользователя
// @Tags         admin
// @Param        id path int true "ID пользователя"
// @Success      200 {object} map[string]string
// @Router       /api/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	if err := h.usersServ.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}

// RegisterDoctorFull godoc
// @Summary      Создать нового врача (аккаунт + профиль)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body RegisterDoctorInput true "Данные нового врача"
// @Router       /api/admin/doctors/create-full [post]
func (h *AdminHandler) RegisterDoctorFull(c *gin.Context) {
	var input RegisterDoctorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Заполните все обязательные поля"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке пароля"})
		return
	}

	newUser := &models.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		MiddleName:   input.MiddleName,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleDoctor,
	}

	if err := h.usersServ.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании аккаунта (возможно, email уже занят)"})
		return
	}

	newDoctor := &models.Doctor{
		UserID:     newUser.UserID,
		Speciality: input.Speciality,
	}

	if err := h.doctorSrv.CreateDoctor(newDoctor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Аккаунт создан, но не удалось назначить специализацию"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Врач успешно зарегистрирован",
		"user_id": newUser.UserID,
	})
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

// GetAllServices godoc
// @Summary      Получение всех сервисов из бд
// @Tags         admin
// @Produce      json
// @Success      200 {array}  dto.ServiceItemDTO
// @Failure      500 {object} map[string]string
// @Router       /api/admin/services [get]
func (h *AdminHandler) GetAllServices(c *gin.Context) {
	serv, err := h.inventorySrv.GetServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список сервисов"})
		return
	}
	c.JSON(http.StatusOK, serv)
}

// UpdateUserRole godoc
// @Summary      Обновить роль пользователя
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path int true "ID пользователя"
// @Param        body body UpdateRoleInput true "Новая роль"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/admin/users/{id}/role [put]
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if err := h.usersServ.UpdateUserRole(id, models.UserRole(input.Role)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Роль обновлена"})
}

// UpdateService godoc
// @Summary      Обновить услугу
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path int true "ID услуги"
// @Param        body body dto.ServiceItemDTO true "Данные услуги"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/admin/services/{id} [put]
func (h *AdminHandler) UpdateService(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var input dto.ServiceItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Данные не верны"})
		return
	}

	if err := h.inventorySrv.UpdateServices(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // ← не было return, 201 слался даже при ошибке
	}
	c.JSON(http.StatusOK, gin.H{"message": "Сервис обновлён"})
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	if err := h.inventorySrv.DeleteService(id); err != nil { // ← была заглушка без реального вызова
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить услугу"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Услуга удалена", "id": id})
}

// GetMeds godoc
// @Summary      Получить все медикаменты
// @Tags         admin
// @Produce      json
// @Success      200 {array}  dto.MedicationItemDTO
// @Router       /api/admin/meds [get]
func (h *AdminHandler) GetMeds(c *gin.Context) {
	meds, err := h.inventorySrv.GetMedications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить медикаменты"})
		return
	}
	c.JSON(http.StatusOK, meds)
}

// CreateMed godoc
// @Summary      Добавить медикамент
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body dto.MedicationItemDTO true "Данные медикамента"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/admin/meds [post]
func (h *AdminHandler) CreateMed(c *gin.Context) {
	var input dto.MedicationItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Данные не верны"})
		return
	}
	if err := h.inventorySrv.CreateMedication(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Медикамент добавлен"})
}

// UpdateMed godoc
// @Summary      Обновить медикамент
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path string true "ID медикамента"
// @Param        body body dto.MedicationItemDTO true "Данные медикамента"
// @Success      200 {object} map[string]string
// @Router       /api/admin/meds/{id} [put]
func (h *AdminHandler) UpdateMed(c *gin.Context) {
	id := c.Param("id")
	var input dto.MedicationItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Данные не верны"})
		return
	}
	if err := h.inventorySrv.UpdateMedication(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Медикамент обновлён"})
}

// DeleteMed godoc
// @Summary      Удалить медикамент
// @Tags         admin
// @Param        id path string true "ID медикамента"
// @Success      200 {object} map[string]string
// @Router       /api/admin/meds/{id} [delete]
func (h *AdminHandler) DeleteMed(c *gin.Context) {
	id := c.Param("id")
	if err := h.inventorySrv.DeleteMedication(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить медикамент"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Медикамент удалён"})
}

// CreateSrv godoc
// @Summary      Добавить сервис
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body dto.ServiceItemDTO true "Данные сервиса"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/admin/services [post]
func (h *AdminHandler) CreateSrv(c *gin.Context) {
	var input dto.ServiceItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Данные не верны"})
		return
	}

	if err := h.inventorySrv.CreateService(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // ← не было return, 201 слался даже при ошибке
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Сервис добавлен"})
}

// CreateDoctor godoc
// @Summary      Создать профиль врача
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body CreateDoctorInput true "Данные врача"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/admin/doctors [post]
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
