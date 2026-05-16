package handler

import (
	"example/project/backend/models"
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	srv service.AppointmentService
}

func NewAppointmentHandler(s service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{srv: s}
}

// Create godoc
// @Summary      Создать запись на приём
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        body body models.Appointment true "Данные записи"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/appointments [post]
func (h *AppointmentHandler) Create(c *gin.Context) {
	var input models.Appointment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные заявки"})
		return
	}

	err := h.srv.CreateAppointment(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Вы успешно записаны на приём"})
}

// GetByOwner godoc
// @Summary      Записи владельца
// @Tags         appointments
// @Produce      json
// @Param        ownerId path int true "ID владельца"
// @Success      200 {array}  models.Appointment
// @Failure      400 {object} map[string]string
// @Router       /api/appointments/owner/{ownerId} [get]
func (h *AppointmentHandler) GetByOwner(c *gin.Context) {
	ownerID, err := strconv.ParseInt(c.Param("ownerId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID владельца"})
		return
	}

	// Мы используем DTO AppointmentRowDTO, чтобы фронтенд сразу видел "🐶 Шарик" и ФИО врача
	list, err := h.srv.GetUpcomingByOwner(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// UpdateStatus godoc
// @Summary      Обновить статус записи
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id   path int true "ID записи"
// @Param        body body object true "Новый статус"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/appointments/{id}/status [put]
func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	appID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input struct {
		Status models.Status `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Статус не указан"})
		return
	}

	err := h.srv.UpdateStatus(appID, input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус обновлен"})
}

// Cancel godoc
// @Summary      Отменить запись
// @Tags         appointments
// @Produce      json
// @Param        id path int true "ID записи"
// @Success      200 {object} map[string]string
// @Router       /api/appointments/{id} [delete]
func (h *AppointmentHandler) Cancel(c *gin.Context) {
	appID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// Отмена — это по сути перевод в статус "rejected"
	err := h.srv.UpdateStatus(appID, models.StatusRejected)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Запись отменена"})
}

func (h *AppointmentHandler) GetBusySlots(c *gin.Context) {
	doctorID, err := strconv.ParseInt(c.Query("doctor_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id"})
		return
	}

	dateStr := c.Query("date") // Ожидаем "2026-05-13"
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}

	slots, err := h.srv.GetOccupiedTimeSlots(doctorID, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slots) // Вернет ["09:00", "11:30"]
}
