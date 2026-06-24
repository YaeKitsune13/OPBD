package handler

import (
	"api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DiagnosisHandler struct {
	svc service.DiagnosisService
}

func NewDiagnosisHandler(svc service.DiagnosisService) *DiagnosisHandler {
	return &DiagnosisHandler{svc}
}

// GetAll godoc
// @Summary      Список диагнозов
// @Description  Возвращает справочник диагнозов для выбора врачом
// @Tags         diagnosis
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  dto.DiagnosisDTO
// @Router       /diagnoses [get]
func (h *DiagnosisHandler) GetAll(c *gin.Context) {
	diagnoses, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, diagnoses)
}
