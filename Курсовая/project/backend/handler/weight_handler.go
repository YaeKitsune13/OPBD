package handler

import (
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WeightHandler struct {
	srv service.PetService
}

func NewWeightHandler(s service.PetService) *WeightHandler {
	return &WeightHandler{srv: s}
}

// GetHistory godoc
// @Summary      История веса питомца
// @Tags         weight
// @Produce      json
// @Param        petId path int true "ID питомца"
// @Success      200 {array}  map[string]interface{}
// @Failure      500 {object} map[string]string
// @Router       /api/weight/pet/{petId} [get]
func (h *WeightHandler) GetHistory(c *gin.Context) {
	petID, _ := strconv.ParseInt(c.Param("petId"), 10, 64)

	data, err := h.srv.GetWeightChartData(petID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// AddRecord godoc
// @Summary      Добавить запись о весе
// @Tags         weight
// @Accept       json
// @Produce      json
// @Param        petId path int    true "ID питомца"
// @Param        body  body object true "Новый вес"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/weight/pet/{petId} [post]
func (h *WeightHandler) AddRecord(c *gin.Context) {
	petID, _ := strconv.ParseInt(c.Param("petId"), 10, 64)

	var input struct {
		Weight float64 `json:"weight"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите корректный вес"})
		return
	}

	// Вызываем метод UpdateWeight, который мы написали в PetService ранее
	// (он и обновит текущий вес, и создаст запись в истории)
	err := h.srv.UpdateWeight(petID, input.Weight, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Вес успешно зафиксирован"})
}
