package handlers

import (
	"strings"

	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/config"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		jwt := strings.Split(auth, " ")[1]
		claims, err := utils.ValidateToken(jwt, config.GetConfig().JWTSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims["uid"].(string))

		c.Next()
	}
}
