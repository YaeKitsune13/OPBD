package repository

import (
	"api/internal/models"

	"gorm.io/gorm"
)

type ShopRepository interface {
	GetProducts() ([]models.Product, error)
	GetCartItems(userID uint) ([]models.CartItem, error)
	AddToCart(item *models.CartItem) error
	UpdateCartQty(itemID uint, qty uint) error
	DeleteCartItem(itemID uint) error
	CreateOrder(order *models.Order) error
	ClearCart(userID uint) error
	GetOrders(userID uint) ([]models.Order, error)
}

type shopRepo struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepo{db}
}

func (r *shopRepo) GetProducts() ([]models.Product, error) {
	products := make([]models.Product, 0)
	err := r.db.Find(&products).Error
	return products, err
}

func (r *shopRepo) GetCartItems(userID uint) ([]models.CartItem, error) {
	items := []models.CartItem{}
	err := r.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (r *shopRepo) AddToCart(item *models.CartItem) error {
	var existing models.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).First(&existing).Error

	if err == nil {
		existing.Quantity += item.Quantity
		return r.db.Save(&existing).Error
	}

	return r.db.Create(item).Error
}

func (r *shopRepo) UpdateCartQty(itemID uint, qty uint) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", qty).Error
}

func (r *shopRepo) DeleteCartItem(itemID uint) error {
	return r.db.Delete(&models.CartItem{}, itemID).Error
}

func (r *shopRepo) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *shopRepo) ClearCart(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}

func (r *shopRepo) GetOrders(userID uint) ([]models.Order, error) {
	orders := []models.Order{}
	err := r.db.Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error
	return orders, err
}
