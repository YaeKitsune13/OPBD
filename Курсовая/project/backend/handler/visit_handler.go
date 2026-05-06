package handler

import (
	"example/project/backend/dto"
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VisitHandler struct {
	srv service.HealthJournalService
}

func NewVisitHandler(s service.HealthJournalService) *VisitHandler {
	return &VisitHandler{srv: s}
}

// SaveVisitCard godoc
// @Summary      Сохранить карту визита
// @Tags         visits
// @Accept       json
// @Produce      json
// @Param        body body dto.ConductVisitDTO true "Данные визита"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/visits [post]
func (h *VisitHandler) SaveVisitCard(c *gin.Context) {
	var input dto.ConductVisitDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные визита"})
		return
	}

	err := h.srv.SaveVisit(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Визит успешно сохранен"})
}

// GetJournal godoc
// @Summary      Журнал здоровья питомца
// @Tags         visits
// @Produce      json
// @Param        petId path int true "ID питомца"
// @Success      200 {array}  dto.ConductVisitDTO
// @Failure      400 {object} map[string]string
// @Router       /api/visits/pet/{petId} [get]
func (h *VisitHandler) GetJournal(c *gin.Context) {
	petID, err := strconv.ParseInt(c.Param("petId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID питомца"})
		return
	}

	journal, err := h.srv.GetPetHistory(petID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, journal)
}

// GetById godoc
// @Summary      Детали визита
// @Tags         visits
// @Produce      json
// @Param        id path int true "ID визита"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/visits/{id} [get]
func (h *VisitHandler) GetById(c *gin.Context) {
	visitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID визита"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"visitId": visitID, "info": "Детальная информация"})
}
