package handler

import (
	"api/internal/dto"
	"api/internal/models"
	"api/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PetHandler struct {
	svc service.PetService
}

func NewPetHandler(svc service.PetService) *PetHandler {
	return &PetHandler{svc}
}

// GetByOwner godoc
// @Summary      Список питомцев владельца
// @Tags         pets
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "ID владельца"
// @Success      200     {array}   dto.PetResponse
// @Router       /pets/owner/{userId} [get]
func (h *PetHandler) GetByOwner(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pets, err := h.svc.GetOwnerPets(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch pets"})
		return
	}

	c.JSON(http.StatusOK, pets)
}

// Create godoc
// @Summary      Добавить питомца
// @Tags         pets
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        ownerId  query     int             true  "ID владельца"
// @Param        request  body      dto.PetRequest  true  "Данные питомца"
// @Success      201      {object}  dto.PetResponse
// @Router       /pets [post]
func (h *PetHandler) Create(c *gin.Context) {
	ownerIdStr := c.Query("ownerId")
	ownerId, _ := strconv.ParseUint(ownerIdStr, 10, 32)

	var req dto.PetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим дату из строки в time.Time
	dob, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	pet := models.Pet{
		OwnerID:   uint(ownerId),
		Name:      req.Name,
		Species:   req.Species,
		Breed:     req.Breed,
		BirthDate: dob,
		Weight:    req.Weight,
		Avatar:    req.Avatar,
	}

	if err := h.svc.AddPet(&pet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save pet"})
		return
	}

	c.JSON(http.StatusCreated, pet)
}

// Update godoc
// @Summary      Редактировать питомца
// @Tags         pets
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        petId    path      int             true  "ID питомца"
// @Param        request  body      dto.PetRequest  true  "Новые данные"
// @Success      200      {object}  dto.PetResponse
// @Router       /pets/{petId} [put]
func (h *PetHandler) Update(c *gin.Context) {
	petIdStr := c.Param("petId")
	petId, _ := strconv.ParseUint(petIdStr, 10, 32)

	var req dto.PetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dob, _ := time.Parse("2006-01-02", req.BirthDate)

	pet := models.Pet{
		Name:      req.Name,
		Species:   req.Species,
		Breed:     req.Breed,
		BirthDate: dob,
		Weight:    req.Weight,
		Avatar:    req.Avatar,
	}
	pet.ID = uint(petId)

	if err := h.svc.UpdatePet(&pet); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pet)
}

// Delete godoc
// @Summary      Удалить питомца
// @Tags         pets
// @Security     ApiKeyAuth
// @Param        petId  path      int  true  "ID питомца"
// @Success      200    {object}  map[string]bool
// @Router       /pets/{petId} [delete]
func (h *PetHandler) Delete(c *gin.Context) {
	petIdStr := c.Param("petId")
	petId, _ := strconv.ParseUint(petIdStr, 10, 32)

	if err := h.svc.RemovePet(uint(petId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
