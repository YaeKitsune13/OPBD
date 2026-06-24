package dto

type UserProfileResponse struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	PetsCount   int    `json:"petsCount"`
	VisitsCount int    `json:"visitsCount"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Phone     string `json:"phone"`
}

type ChangePasswordRequest struct {
	Current string `json:"current" binding:"required"`
	Next    string `json:"next" binding:"required,min=8"`
}
