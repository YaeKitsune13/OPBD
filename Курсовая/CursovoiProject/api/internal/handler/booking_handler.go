package handler

import (
	"api/internal/dto"
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	svc service.BookingService
}

func NewBookingHandler(svc service.BookingService) *BookingHandler {
	return &BookingHandler{svc}
}

// GetInit godoc
// @Summary      Инициализация формы записи
// @Description  Получает списки питомцев, врачей и услуг
// @Tags         booking
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "ID пользователя"
// @Success      200     {object}  dto.BookingInitResponse
// @Router       /book/init/{userId} [get]
func (h *BookingHandler) GetInit(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	data, _ := h.svc.GetInitData(uint(uid))
	c.JSON(http.StatusOK, data)
}

// GetBusySlots godoc
// @Summary      Получение занятых слотов
// @Description  Возвращает список времени, на которое уже есть запись у врача
// @Tags         booking
// @Security     ApiKeyAuth
// @Produce      json
// @Param        doctor_id  query     int     true  "ID врача"
// @Param        date       query     string  true  "Дата (YYYY-MM-DD)"
// @Success      200        {array}   string
// @Router       /appointments/busy-slots [get]
func (h *BookingHandler) GetBusySlots(c *gin.Context) {
	did, _ := strconv.ParseUint(c.Query("doctor_id"), 10, 32)
	date := c.Query("date") // "2024-05-20"
	slots, _ := h.svc.GetBusySlots(uint(did), date)
	c.JSON(http.StatusOK, slots)
}

// Create godoc
// @Summary      Создание записи на прием
// @Description  Бронирует визит к врачу
// @Tags         booking
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateAppointmentRequest  true  "Данные записи"
// @Success      201      {object}  map[string]bool
// @Router       /appointments [post]
func (h *BookingHandler) Create(c *gin.Context) {
	var req dto.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID := c.GetUint("userID")

	if err := h.svc.CreateAppointment(clientID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true})
}

// GetHistory godoc
// @Summary      История визитов клиента
// @Tags         booking
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "ID пользователя"
// @Success      200     {array}   dto.RecentAppointmentDTO
// @Router       /appointments/client/{userId} [get]
func (h *BookingHandler) GetHistory(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	history, err := h.svc.GetClientHistory(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
