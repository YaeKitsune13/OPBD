package handler

import (
	"example/project/backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchSrv service.SearchService
}

func NewSearchHandler(s service.SearchService) *SearchHandler {
	return &SearchHandler{searchSrv: s}
}

// Search godoc
// @Summary      Поиск пациентов
// @Tags         search
// @Produce      json
// @Param        q query string true "Кличка питомца или ФИО владельца"
// @Success      200 {array} dto.PatientSearchResultDTO
// @Router       /api/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")

	results, err := h.searchSrv.SearchPatients(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска"})
		return
	}

	c.JSON(http.StatusOK, results)
}
