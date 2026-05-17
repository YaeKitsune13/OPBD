package models

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Структура данных внутри токена (Claims)
type MyCustomClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GetJWTSecret возвращает ключ из переменной окружения.
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("============================================")
		log.Println(".env not access jwt token")
		log.Println("============================================")
		return []byte("fallback_secret_key_for_dev")
	}
	return []byte(secret)
}
