package handlers

import (
	"strconv"

	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/models"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

// AddItem adds an item to the shopping cart
func (h *CartHandler) AddItem(c *gin.Context) {
	uid := c.MustGet("userID")

	var item models.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	response, err := h.service.AddItem(c, uid.(string), &item)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

// SetQuantity sets the quantity of an item in the shopping cart
func (h *CartHandler) SetQuantity(c *gin.Context) {
	uid := c.MustGet("userID")

	itemIDStr := c.Param("itemID")

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	type setQuantityRequest struct {
		Quantity int `json:"quantity"`
	}

	var req setQuantityRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	response, err := h.service.SetQuantity(c, uid.(string), itemID, req.Quantity)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

// RemoveItem removes an item from the shopping cart
func (h *CartHandler) RemoveItem(c *gin.Context) {
	uid := c.MustGet("userID")

	itemIDStr := c.Param("itemID")

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	response, err := h.service.RemoveItem(c, uid.(string), itemID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

// GetItems returns all items in the shopping cart
func (h *CartHandler) GetItems(c *gin.Context) {
	uid := c.MustGet("userID")

	response, err := h.service.GetItems(c, uid.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

// DeleteCart deletes the shopping cart
func (h *CartHandler) DeleteCart(c *gin.Context) {
	uid := c.MustGet("userID")

	err := h.service.DeleteCart(c, uid.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Cart deleted"})
}
