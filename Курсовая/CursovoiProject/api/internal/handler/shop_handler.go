package handler

import (
	"api/internal/dto"
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	svc service.ShopService
}

func NewShopHandler(svc service.ShopService) *ShopHandler {
	return &ShopHandler{svc}
}

// GetProducts godoc
// @Summary      Каталог товаров
// @Tags         shop
// @Produce      json
// @Success      200  {array}  dto.ProductResponse
// @Router       /medications [get]
func (h *ShopHandler) AddToCart(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		ProductID uint `json:"productId" binding:"required"`
		Quantity  uint `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные товара: " + err.Error()})
		return
	}

	if err := h.svc.AddItemToCart(uint(uid), req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить в корзину"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetCart godoc
// @Summary      Получить корзину
// @Tags         shop
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "ID пользователя"
// @Success      200     {array}   dto.CartItemResponse
// @Router       /cart/{userId} [get]2
func (h *ShopHandler) GetProducts(c *gin.Context) {
	products, err := h.svc.GetCatalog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// AddToCart godoc
// @Summary      Добавить в корзину
// @Tags         shop
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        userId   path      int  true  "ID пользователя"
// @Param        request  body      object{productId=int,quantity=int}  true  "ID товара и количество"
// @Success      200      {object}  map[string]bool
// @Router       /cart/{userId} [post]
func (h *ShopHandler) GetCart(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	cart, err := h.svc.GetUserCart(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// UpdateCart godoc
// @Summary      Изменить количество в корзине
// @Tags         shop
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        itemId   path      int  true  "ID элемента корзины"
// @Param        request  body      object{quantity=int}  true  "Новое количество"
// @Success      200      {object}  map[string]bool
// @Router       /cart/{itemId} [put]
func (h *ShopHandler) UpdateCart(c *gin.Context) {
	itemID, _ := strconv.ParseUint(c.Param("itemId"), 10, 32)
	var req struct {
		Quantity uint `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.svc.ChangeCartQty(uint(itemID), req.Quantity)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteCart godoc
// @Summary      Удалить из корзины
// @Tags         shop
// @Security     ApiKeyAuth
// @Param        itemId  path      int  true  "ID элемента корзины"
// @Success      200     {object}  map[string]bool
// @Router       /cart/{itemId} [delete]
func (h *ShopHandler) DeleteCart(c *gin.Context) {
	itemID, _ := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err := h.svc.RemoveFromCart(uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Checkout godoc
// @Summary      Оформить заказ
// @Tags         shop
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.OrderRequest  true  "Данные заказа"
// @Success      200      {object}  map[string]bool
// @Router       /orders [post]
func (h *ShopHandler) Checkout(c *gin.Context) {
	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Checkout(req.UserID, req.TotalAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetOrders godoc
// @Summary      История заказов
// @Description  Получает список всех заказов пользователя с деталями товаров для чека
// @Tags         shop
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "ID пользователя"
// @Success      200     {array}   dto.OrderResponse
// @Router       /orders/{userId} [get]
func (h *ShopHandler) GetOrders(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	orders, err := h.svc.GetOrders(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
