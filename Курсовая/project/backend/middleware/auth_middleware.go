package middleware

import (
	"net/http"
	"strings"

	"example/project/backend/models" // замени на свой путь

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет наличие и валидность JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Достаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Заголовок Authorization отсутствует"})
			c.Abort() // Останавливаем запрос
			return
		}

		// 2. Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат заголовка (Bearer <token>)"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Парсим и валидируем токен
		claims := &models.MyCustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return models.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный или просроченный токен"})
			c.Abort()
			return
		}

		// 4. Сохраняем данные пользователя в контекст Gin,
		// чтобы хендлеры могли их использовать (например, узнать ID того, кто делает запрос)
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next() // Пропускаем запрос дальше
	}
}

// RoleMiddleware проверяет, есть ли у пользователя права (например, только для admin)
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Роль не определена"})
			c.Abort()
			return
		}

		userRole := role.(string)
		isAllowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "У вас недостаточно прав"})
			c.Abort()
			return
		}

		c.Next()
	}
}
