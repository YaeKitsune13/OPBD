package handler

import (
	"api/internal/dto"
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	svc service.DoctorService
}

func NewDoctorHandler(svc service.DoctorService) *DoctorHandler {
	return &DoctorHandler{svc}
}

// GetSchedule godoc
// @Summary      Расписание врача
// @Description  Список всех записей к конкретному врачу
// @Tags         doctor
// @Security     ApiKeyAuth
// @Produce      json
// @Param        doctor_id  query     int  true  "ID врача"
// @Success      200        {array}   dto.DoctorScheduleDTO
// @Router       /doctor/schedule [get]
func (h *DoctorHandler) GetSchedule(c *gin.Context) {
	did, _ := strconv.ParseUint(c.Query("doctor_id"), 10, 32)
	c.JSON(http.StatusOK, h.svc.GetSchedule(uint(did)))
}

// CompleteVisit godoc
// @Summary      Завершение приема
// @Description  Закрывает визит и сохраняет медицинский протокол (диагноз, вес, лечение)
// @Tags         doctor
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "ID записи"
// @Param        request  body      dto.CompleteVisitRequest  true  "Данные осмотра"
// @Success      200      {object}  map[string]bool
// @Router       /appointments/{id}/complete [patch]
func (h *DoctorHandler) CompleteVisit(c *gin.Context) {
	aid, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req dto.CompleteVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.svc.CompleteVisit(uint(aid), req)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPatients godoc
// @Summary      Search patients
// @Description  Search for clients by name or phone
// @Tags         doctor
// @Security     ApiKeyAuth
// @Produce      json
// @Param        search  query     string  false  "Search query"
// @Success      200     {array}   dto.PatientDTO
// @Router       /doctor/patients [get]
func (h *DoctorHandler) GetPatients(c *gin.Context) {
	search := c.Query("search")
	c.JSON(http.StatusOK, h.svc.GetPatients(search))
}

// GetHistory godoc
// @Summary      Медицинская история клиента
// @Description  Получает историю всех визитов всех питомцев конкретного клиента
// @Tags         doctor
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      int  true  "ID клиента"
// @Success      200  {array}   dto.PetHistoryDTO
// @Router       /doctor/patients/{id}/history [get]
func (h *DoctorHandler) GetHistory(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	c.JSON(http.StatusOK, h.svc.GetMedicalHistory(uint(cid)))
}

// UpdateStatus godoc
// @Summary      Изменить статус записи
// @Description  Принять (confirmed) или отклонить (rejected) запись
// @Tags         doctor
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "ID записи"
// @Param        request  body      object{status=string}  true  "Новый статус"
// @Success      200      {object}  map[string]bool
// @Router       /appointments/{id}/status [patch]
func (h *DoctorHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID записи"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле status обязательно"})
		return
	}

	if err := h.svc.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить статус"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
