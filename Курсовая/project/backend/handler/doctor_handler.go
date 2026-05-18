package handler

import (
	"example/project/backend/models"
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorSrv service.DoctorService
	appSrv    service.AppointmentService
}

func NewDoctorHandler(ds service.DoctorService, as service.AppointmentService) *DoctorHandler {
	return &DoctorHandler{
		doctorSrv: ds,
		appSrv:    as,
	}
}

// GetSchedule godoc
// @Summary      Расписание врача на сегодня
// @Tags         doctors
// @Produce      json
// @Param        id path int true "ID врача"
// @Success      200 {array}  models.Appointment
// @Failure      400 {object} map[string]string
// @Router       /api/doctors/{id}/schedule [get]
func (h *DoctorHandler) GetSchedule(c *gin.Context) {
	doctorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID врача"})
		return
	}

	schedule, err := h.appSrv.GetDoctorTodaySchedule(doctorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetBySpecialty godoc
// @Summary      Список врачей (все или по специализации)
// @Tags         doctors
// @Produce      json
// @Param        specialty query string false "Специализация"
// @Success      200 {array}  models.Doctor
// @Failure      500 {object} map[string]string
// @Router       /api/doctors [get]
func (h *DoctorHandler) GetBySpecialty(c *gin.Context) {
	specialty := c.Query("specialty") // Например: /api/doctors?specialty=Хирург

	var doctors []models.Doctor
	var err error

	if specialty != "" {
		doctors, err = h.doctorSrv.GetBySpecialty(specialty)
	} else {
		doctors, err = h.doctorSrv.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

// GetMe godoc
// @Summary      Получить профиль текущего доктора
// @Tags         doctors
// @Produce      json
// @Success      200 {object} models.Doctor
// @Router       /api/doctors/me [get]
func (h *DoctorHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")
	doctor, err := h.doctorSrv.GetByUserID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль доктора не найден"})
		return
	}
	c.JSON(http.StatusOK, doctor)
}

// GetTodaySchedule godoc
// @Summary      Расписание доктора на сегодня
// @Tags         doctors
// @Produce      json
// @Success      200 {array} dto.TodayScheduleDTO
// @Router       /api/doctors/me/schedule [get]
func (h *DoctorHandler) GetTodaySchedule(c *gin.Context) {
	userID, _ := c.Get("userID")
	doctor, err := h.doctorSrv.GetByUserID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль доктора не найден"})
		return
	}

	schedule, err := h.appSrv.GetDoctorTodaySchedule(doctor.DoctorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки расписания"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}
