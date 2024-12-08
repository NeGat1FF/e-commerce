package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/NeGat1FF/e-commerce/search-service/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(s *service.SearchService) *SearchHandler {
	return &SearchHandler{
		service: s,
	}
}

func (h *SearchHandler) SearchProducts(c *gin.Context) {
	rawQuery := c.Request.URL.RawQuery

	// Split the query string into key-value pairs
	queryParams := strings.Split(rawQuery, "&")
	params := make(map[string]string)
	var priceMin, priceMax int
	var sortOrder string
	for _, param := range queryParams {
		pair := strings.Split(param, "=")
		key := pair[0]
		value := pair[1]
		switch key {
		case "sort":
			sortOrder = value
		case "price_min":
			priceMin, _ = strconv.Atoi(value)
		case "price_max":
			priceMax, _ = strconv.Atoi(value)
		default:
			params[key] = value
		}
	}

	products, err := h.service.SearchProducts(c, params, sortOrder, priceMin, priceMax, 1, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}
