package dto

import "example/project/backend/models"

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Phone     string `json:"phone"`
}

type AuthResponse struct {
	Token    string          `json:"token"`
	Role     models.UserRole `json:"role"` // "admin", "doctor", "client"
	UserName string          `json:"userName"`
	UserID   int64           `json:"userId"`
}
