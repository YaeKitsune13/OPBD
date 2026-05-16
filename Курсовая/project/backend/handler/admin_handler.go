package handler

import (
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

	// 1. Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке пароля"})
		return
	}

	// 2. Подготавливаем модель пользователя
	newUser := &models.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		MiddleName:   input.MiddleName,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleDoctor, // Убедитесь, что константа RoleDoctor определена в моделях
	}

	// 3. Сохраняем пользователя в БД через сервис
	if err := h.usersServ.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании аккаунта (возможно, email уже занят)"})
		return
	}

	// 4. Создаем запись в таблице врачей, используя полученный UserID
	newDoctor := &models.Doctor{
		UserID:     newUser.UserID,
		Speciality: input.Speciality,
	}

	if err := h.doctorSrv.CreateDoctor(newDoctor); err != nil {
		// Опционально: здесь можно удалить созданного пользователя, если создание профиля врача провалилось (транзакция)
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
