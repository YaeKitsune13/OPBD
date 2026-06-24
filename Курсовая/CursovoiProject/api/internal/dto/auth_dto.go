package dto

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	Phone     string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserId   uint   `json:"userId"`
	Role     string `json:"role"`
	UserName string `json:"userName"`
	LastName string `json:"lastName"`
	Phone    string `json:"phone"`
}
