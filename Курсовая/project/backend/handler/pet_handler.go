package handler

import (
	"example/project/backend/dto"
	"example/project/backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PetHandler struct {
	srv service.PetService
}

func NewPetHandler(s service.PetService) *PetHandler {
	return &PetHandler{srv: s}
}

// GetByOwner godoc
// @Summary      Список питомцев владельца
// @Tags         pets
// @Produce      json
// @Param        ownerId path int true "ID владельца"
// @Success      200 {array}  dto.PetCardDTO
// @Failure      400 {object} map[string]string
// @Router       /api/pets/owner/{ownerId} [get]
func (h *PetHandler) GetByOwner(c *gin.Context) {
	// Извлекаем :ownerId из URL
	ownerID, err := strconv.ParseInt(c.Param("ownerId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID владельца"})
		return
	}

	pets, err := h.srv.GetOwnerPets(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pets)
}

// AddPet godoc
// @Summary      Добавить питомца
// @Tags         pets
// @Accept       json
// @Produce      json
// @Param        ownerId query int    true "ID владельца"
// @Param        body    body  dto.PetCardDTO true "Данные питомца"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/pets [post]
func (h *PetHandler) AddPet(c *gin.Context) {
	var input dto.PetCardDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		// Теперь эта ошибка вылетит только если JSON совсем "битый"
		c.JSON(http.StatusBadRequest, gin.H{"error": "Проверьте правильность заполнения полей"})
		return
	}

	ownerID, _ := strconv.ParseInt(c.Query("ownerId"), 10, 64)

	err := h.srv.AddPet(ownerID, input)
	if err != nil {
		// Возвращаем текст ошибки из сервиса (например, "дата рождения не может быть в будущем")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Питомец добавлен"})
}

// UpdatePet godoc
// @Summary      Обновить данные питомца
// @Tags         pets
// @Accept       json
// @Produce      json
// @Param        petId path int          true "ID питомца"
// @Param        body  body dto.PetCardDTO true "Новые данные"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/pets/{petId} [put]
func (h *PetHandler) UpdatePet(c *gin.Context) {
	var input dto.PetCardDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	petID, _ := strconv.ParseInt(c.Param("petId"), 10, 64)
	input.PetId = petID

	err := h.srv.UpdatePet(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Данные обновлены"})
}

// DeletePet godoc
// @Summary      Удалить питомца
// @Tags         pets
// @Produce      json
// @Param        petId path int true "ID питомца"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/pets/{petId} [delete]
func (h *PetHandler) DeletePet(c *gin.Context) {
	petID, err := strconv.ParseInt(c.Param("petId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID питомца"})
		return
	}

	err = h.srv.DeletePet(petID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Питомец удален"})
}
