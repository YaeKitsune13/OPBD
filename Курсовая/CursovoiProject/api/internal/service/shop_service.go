package service

import (
	"api/internal/dto"
	"api/internal/models"
	"api/internal/repository"
	"errors"
)

type ShopService interface {
	GetCatalog() ([]models.Product, error)
	GetUserCart(userID uint) ([]dto.CartItemResponse, error)
	AddItemToCart(userID uint, productID uint, qty uint) error
	ChangeCartQty(itemID uint, qty uint) error
	RemoveFromCart(itemID uint) error
	Checkout(userID uint, total uint) error
	GetOrders(userID uint) ([]dto.OrderResponse, error)
}

type shopService struct {
	repo repository.ShopRepository
}

func NewShopService(repo repository.ShopRepository) ShopService {
	return &shopService{repo}
}

func (s *shopService) GetCatalog() ([]models.Product, error) {
	return s.repo.GetProducts()
}

func (s *shopService) GetUserCart(userID uint) ([]dto.CartItemResponse, error) {
	items, err := s.repo.GetCartItems(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.CartItemResponse
	for _, item := range items {
		res = append(res, dto.CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Price:     item.Product.Price,
			Quantity:  item.Quantity,
		})
	}
	return res, nil
}

func (s *shopService) AddItemToCart(userID uint, productID uint, qty uint) error {
	return s.repo.AddToCart(&models.CartItem{UserID: userID, ProductID: productID, Quantity: qty})
}

func (s *shopService) ChangeCartQty(itemID uint, qty uint) error {
	return s.repo.UpdateCartQty(itemID, qty)
}

func (s *shopService) RemoveFromCart(itemID uint) error {
	return s.repo.DeleteCartItem(itemID)
}

func (s *shopService) Checkout(userID uint, total uint) error {
	cartItems, err := s.repo.GetCartItems(userID)
	if err != nil || len(cartItems) == 0 {
		return errors.New("корзина пуста")
	}

	order := models.Order{
		UserID:      userID,
		TotalAmount: total,
		Status:      models.OrderPaid,
		Items:       make([]models.OrderItem, 0),
	}

	for _, ci := range cartItems {
		order.Items = append(order.Items, models.OrderItem{
			ProductID: ci.ProductID,
			Quantity:  ci.Quantity,
			Price:     ci.Product.Price,
		})
	}

	if err := s.repo.CreateOrder(&order); err != nil {
		return err
	}

	return s.repo.ClearCart(userID)
}

func (s *shopService) GetOrders(userID uint) ([]dto.OrderResponse, error) {
	orders, err := s.repo.GetOrders(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.OrderResponse
	for _, o := range orders {
		var items []dto.OrderItemResponse
		for _, item := range o.Items {
			items = append(items, dto.OrderItemResponse{
				ProductName: item.Product.Name,
				Price:       item.Price,
				Quantity:    item.Quantity,
			})
		}

		res = append(res, dto.OrderResponse{
			ID:          o.ID,
			Date:        o.CreatedAt.Format("02.01.2006 15:04"),
			TotalAmount: o.TotalAmount,
			Status:      string(o.Status),
			Items:       items,
		})
	}
	return res, nil
}
