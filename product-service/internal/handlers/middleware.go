package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/NeGat1FF/e-commerce/product-service/internal/models"
	"github.com/NeGat1FF/e-commerce/product-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Logging middleware logs the incoming requests
func Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		execTime := time.Since(start)
		clientIP := ctx.ClientIP()

		logger.Logger.Info("Request", zap.String("client_ip", clientIP), zap.String("method", ctx.Request.Method), zap.String("path", ctx.Request.URL.Path), zap.Int("status", ctx.Writer.Status()), zap.Duration("exec_time", execTime))
	}
}

// ErrorHandling middleware handles the errors
func ErrorHandling() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				logger.Logger.Error("Error occurred", zap.Error(err.Err))
			}
		}
	}
}

// Auth middleware reads user role from JWT token and sets it to the context
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.JSON(401, gin.H{
				"error": "authorization header is required",
			})
			ctx.Abort()
			return
		}

		tokenString := strings.Split(header, " ")[1]

		claims := jwt.MapClaims{}

		_, _, err := jwt.NewParser().ParseUnverified(tokenString, &claims)
		if err != nil {
			ctx.JSON(401, gin.H{
				"error": "failed to parse token",
			})
			ctx.Abort()
			return
		}

		role, ok := claims["role"]
		if !ok {
			ctx.JSON(401, gin.H{
				"error": "role not found in token",
			})
			ctx.Abort()
			return
		}

		if role != "admin" {
			ctx.JSON(401, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func ValidateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product models.Product
		var productMap map[string]interface{}

		if ctx.Request.Method == "POST" {
			// Bind to model for create operation
			if err := ctx.ShouldBindJSON(&product); err != nil {
				fmt.Println(err)
				ctx.JSON(400, gin.H{
					"error": err.Error(),
				})
				ctx.Abort()
				return
			}

			if product.Price <= 0 {
				ctx.JSON(400, gin.H{
					"error": "price must be greater than 0",
				})
				ctx.Abort()
				return
			}

			if product.Quantity < 0 {
				ctx.JSON(400, gin.H{
					"error": "quantity must be greater than or equal to 0",
				})
				ctx.Abort()
				return
			}

			ctx.Set("product", product)
		} else if ctx.Request.Method == "PUT" {
			// Bind to map for update operation
			if err := ctx.ShouldBindJSON(&productMap); err != nil {
				ctx.JSON(400, gin.H{
					"error": err.Error(),
				})
				ctx.Abort()
				return
			}

			if price, ok := productMap["price"].(float64); ok && price <= 0 {
				ctx.JSON(400, gin.H{
					"error": "price must be greater than 0",
				})
				ctx.Abort()
				return
			}

			if _, ok := productMap["quantity"]; ok {
				ctx.JSON(400, gin.H{
					"error": "quantity cannot be updated, use /add-stock or /remove-stock endpoints",
				})
				ctx.Abort()
				return
			}

			ctx.Set("updateFields", productMap)
		}

		ctx.Next()
	}
}
