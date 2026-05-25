package handler

import (
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	svc    service.StatsService
	petSvc service.PetService
}

func NewStatsHandler(svc service.StatsService, petSvc service.PetService) *StatsHandler {
	return &StatsHandler{svc, petSvc}
}

// GetUserPets godoc
// @Summary      Список питомцев пользователя (краткий)
// @Tags         stats
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "ID пользователя"
// @Success      200     {array}   object{id=int,name=string}
// @Router       /stats/pets/{userId} [get]
func (h *StatsHandler) GetUserPets(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	pets, err := h.petSvc.GetOwnerPets(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем упрощенный список для селектора
	type brief struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	var res []brief
	for _, p := range pets {
		res = append(res, brief{ID: p.ID, Name: p.Name})
	}
	c.JSON(http.StatusOK, res)
}

// GetWeightData godoc
// @Summary      Данные веса для графика
// @Tags         stats
// @Security     ApiKeyAuth
// @Produce      json
// @Param        petId  path      int  true  "ID питомца"
// @Success      200    {array}   dto.WeightDataDTO
// @Router       /stats/weight/{petId} [get]
func (h *StatsHandler) GetWeightData(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("petId"), 10, 32)
	data, err := h.svc.GetPetWeightStats(uint(pid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, data)
}
