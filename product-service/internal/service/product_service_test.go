package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NeGat1FF/e-commerce/product-service/internal/models"
	"github.com/NeGat1FF/e-commerce/product-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/product-service/internal/service"
	"github.com/NeGat1FF/e-commerce/product-service/mocks"
	"github.com/NeGat1FF/e-commerce/product-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		name string
		models.Product
		setupMocks    func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue)
		expectedError error
	}{
		{
			name: "Create product success",
			Product: models.Product{
				Name:     "Test Product",
				Price:    100,
				Quantity: 10,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("CreateProduct", mock.Anything, mock.Anything).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.created", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Create product failed in repository",
			Product: models.Product{
				Name:     "Test Product",
				Price:    100,
				Quantity: 10,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("CreateProduct", mock.Anything, mock.Anything).Return(errors.New("failed to add product to database"))
			},
			expectedError: service.ErrFailedToCreateProduct,
		},
		{
			name: "Faield to publish message to message queue",
			Product: models.Product{
				Name:     "Test Product",
				Price:    100,
				Quantity: 10,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("CreateProduct", mock.Anything, mock.Anything).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.created", mock.Anything).Return(errors.New("failed to publish message"))
			},
			expectedError: nil,
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}
			cache := &mocks.Cache{}
			messageQueue := &mocks.MessageQueue{}

			tc.setupMocks(repository, cache, messageQueue)

			productService := service.NewProductService(repository, messageQueue, cache, "test-exchange")

			err := productService.CreateProduct(context.Background(), tc.Product)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			}

			time.Sleep(100 * time.Millisecond) // wait for goroutine to finish

			repository.AssertExpectations(t)
			cache.AssertExpectations(t)
			messageQueue.AssertExpectations(t)
		})
	}

}

func TestUpdateProduct(t *testing.T) {
	testCases := []struct {
		name          string
		productID     int64
		updateFields  map[string]interface{}
		setupMocks    func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue)
		expectedError error
	}{
		{
			name:      "Update product success",
			productID: 1,
			updateFields: map[string]interface{}{
				"price": 200,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("UpdateProduct", mock.Anything, mock.AnythingOfType("int64"), mock.Anything).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.updated", mock.Anything).Return(nil)
				c.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Product not found",
			productID: 1,
			updateFields: map[string]interface{}{
				"price": 200,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("UpdateProduct", mock.Anything, mock.AnythingOfType("int64"), mock.Anything).Return(repository.ErrProductNotFound)
			},
			expectedError: service.ErrProductNotFound,
		},
		{
			name:      "Update product failed in repository",
			productID: 1,
			updateFields: map[string]interface{}{
				"price": 200,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("UpdateProduct", mock.Anything, mock.AnythingOfType("int64"), mock.Anything).Return(errors.New("failed to update product in database"))
			},
			expectedError: service.ErrFailedToUpdateProduct,
		},
		{
			name:      "Failed to publish message to message queue",
			productID: 1,
			updateFields: map[string]interface{}{
				"price": 200,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("UpdateProduct", mock.Anything, mock.AnythingOfType("int64"), mock.Anything).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.updated", mock.Anything).Return(errors.New("failed to publish message"))
				c.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Failed to delete product from cache",
			productID: 1,
			updateFields: map[string]interface{}{
				"price": 200,
			},
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("UpdateProduct", mock.Anything, mock.AnythingOfType("int64"), mock.Anything).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.updated", mock.Anything).Return(nil)
				c.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("failed to delete product from cache"))
			},
			expectedError: nil,
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}
			cache := &mocks.Cache{}
			messageQueue := &mocks.MessageQueue{}

			tc.setupMocks(repository, cache, messageQueue)

			productService := service.NewProductService(repository, messageQueue, cache, "test-exchange")

			err := productService.UpdateProduct(context.Background(), 1, tc.updateFields)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			}

			time.Sleep(100 * time.Millisecond) // wait for goroutine to finish

			repository.AssertExpectations(t)
			cache.AssertExpectations(t)
			messageQueue.AssertExpectations(t)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	testCases := []struct {
		name          string
		productID     int64
		setupMocks    func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue)
		expectedError error
	}{
		{
			name:      "Delete product success",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("DeleteProduct", mock.Anything, mock.AnythingOfType("int64")).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.deleted", mock.Anything).Return(nil)
				c.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Product not found",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("DeleteProduct", mock.Anything, mock.AnythingOfType("int64")).Return(repository.ErrProductNotFound)
			},
			expectedError: service.ErrProductNotFound,
		},
		{
			name:      "Delete product failed in repository",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("DeleteProduct", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("failed to delete product in database"))
			},
			expectedError: service.ErrFailedToDeleteProduct,
		},
		{
			name:      "Failed to publish message to message queue",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("DeleteProduct", mock.Anything, mock.AnythingOfType("int64")).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.deleted", mock.Anything).Return(errors.New("failed to publish message"))
				c.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Failed to delete product from cache",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache, m *mocks.MessageQueue) {
				r.On("DeleteProduct", mock.Anything, mock.AnythingOfType("int64")).Return(nil)
				m.On("PublishMessage", mock.Anything, mock.Anything, "product.deleted", mock.Anything).Return(nil)
				c.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("failed to delete product from cache"))
			},
			expectedError: nil,
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}
			cache := &mocks.Cache{}
			messageQueue := &mocks.MessageQueue{}

			tc.setupMocks(repository, cache, messageQueue)

			productService := service.NewProductService(repository, messageQueue, cache, "test-exchange")

			err := productService.DeleteProduct(context.Background(), 1)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			}

			time.Sleep(100 * time.Millisecond) // wait for goroutine to finish

			repository.AssertExpectations(t)
			cache.AssertExpectations(t)
			messageQueue.AssertExpectations(t)
		})
	}
}

func TestGetProductByID(t *testing.T) {
	testCases := []struct {
		name            string
		productID       int64
		setupMocks      func(r *mocks.ProductRepository, c *mocks.Cache)
		expectedProduct models.UserProduct
		expectedError   error
	}{
		{
			name:      "Get product by id success",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache) {
				r.On("GetProductByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.UserProduct{
					ID:    1,
					Name:  "Test Product",
					Price: 100,
				}, nil)
				c.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("failed to get product from cache"))
				c.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("models.UserProduct")).Return(nil)
			},
			expectedProduct: models.UserProduct{
				ID:    1,
				Name:  "Test Product",
				Price: 100,
			},
			expectedError: nil,
		},
		{
			name:      "Get product by id from cache",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache) {
				c.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.UserProduct)
					product.ID = 1
					product.Name = "Test Product"
					product.Price = 100
				})
			},
			expectedProduct: models.UserProduct{
				ID:    1,
				Name:  "Test Product",
				Price: 100,
			},
		},
		{
			name:      "Product not found",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache) {
				r.On("GetProductByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.UserProduct{}, repository.ErrProductNotFound)
				c.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("failed to get product from cache"))
			},
			expectedProduct: models.UserProduct{},
			expectedError:   service.ErrProductNotFound,
		},
		{
			name:      "Failed to get product by id",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache) {
				r.On("GetProductByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.UserProduct{}, errors.New("failed to get product by id"))
				c.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("failed to get product from cache"))
			},
			expectedProduct: models.UserProduct{},
			expectedError:   service.ErrFailedToGetProductByID,
		},
		{
			name:      "Failed to set product to cache",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache) {
				r.On("GetProductByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.UserProduct{
					ID:    1,
					Name:  "Test Product",
					Price: 100,
				}, nil)
				c.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("failed to get product from cache"))
				c.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("models.UserProduct")).Return(errors.New("failed to set product to cache"))
			},
			expectedProduct: models.UserProduct{
				ID:    1,
				Name:  "Test Product",
				Price: 100,
			},
			expectedError: nil,
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}
			cache := &mocks.Cache{}

			tc.setupMocks(repository, cache)

			productService := service.NewProductService(repository, nil, cache, "")

			product, err := productService.GetProductByID(context.Background(), tc.productID)

			assert.Equal(t, tc.expectedProduct, product)
			assert.Equal(t, tc.expectedError, err)

			time.Sleep(100 * time.Millisecond) // wait for goroutine to finish

			repository.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}

func TestGetProductByCategory(t *testing.T) {
	testCases := []struct {
		name             string
		category         string
		page             int
		limit            int
		setupMocks       func(r *mocks.ProductRepository, c *mocks.Cache)
		expectedProducts []models.UserProduct
		expectedError    error
	}{
		{
			name:     "Get products by category success",
			category: "test",
			page:     1,
			limit:    10,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache) {
				r.On("GetProductsByCategory", mock.Anything, "test", 1, 10).Return([]models.UserProduct{
					{
						ID:    1,
						Name:  "Test Product",
						Price: 100,
					},
				}, nil)
			},
			expectedProducts: []models.UserProduct{
				{
					ID:    1,
					Name:  "Test Product",
					Price: 100,
				},
			},
			expectedError: nil,
		},
		{
			name:     "Failed to get products by category",
			category: "test",
			page:     1,
			limit:    10,
			setupMocks: func(r *mocks.ProductRepository, c *mocks.Cache) {
				r.On("GetProductsByCategory", mock.Anything, "test", 1, 10).Return(nil, errors.New("failed to get products by category"))
			},
			expectedProducts: nil,
			expectedError:    errors.New("failed to get products by category"),
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}
			cache := &mocks.Cache{}

			tc.setupMocks(repository, cache)

			productService := service.NewProductService(repository, nil, cache, "")

			products, err := productService.GetProductsByCategory(context.Background(), tc.category, tc.page, tc.limit)

			assert.Equal(t, tc.expectedProducts, products)
			assert.Equal(t, tc.expectedError, err)

			time.Sleep(100 * time.Millisecond) // wait for goroutine to finish

			repository.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}

func TestGetStock(t *testing.T) {
	testCases := []struct {
		name          string
		productID     int64
		setupMocks    func(r *mocks.ProductRepository)
		expectedStock int64
		expectedError error
	}{
		{
			name:      "Get stock success",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository) {
				r.On("GetStock", mock.Anything, int64(1)).Return(int64(10), nil)
			},
			expectedStock: int64(10),
			expectedError: nil,
		},
		{
			name:      "Failed to get stock",
			productID: 1,
			setupMocks: func(r *mocks.ProductRepository) {
				r.On("GetStock", mock.Anything, int64(1)).Return(int64(0), errors.New("failed to get stock"))
			},
			expectedStock: int64(0),
			expectedError: errors.New("failed to get stock"),
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}

			tc.setupMocks(repository)

			productService := service.NewProductService(repository, nil, nil, "")

			stock, err := productService.GetStock(context.Background(), tc.productID)

			assert.Equal(t, tc.expectedStock, stock)
			assert.Equal(t, tc.expectedError, err)

			repository.AssertExpectations(t)
		})
	}
}

func TestAddStock(t *testing.T) {
	testCases := []struct {
		name          string
		productID     int64
		quantity      int64
		setupMocks    func(r *mocks.ProductRepository)
		expectedError error
	}{
		{
			name:      "Add stock success",
			productID: 1,
			quantity:  10,
			setupMocks: func(r *mocks.ProductRepository) {
				r.On("AddStock", mock.Anything, int64(1), int64(10)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Failed to add stock",
			productID: 1,
			quantity:  10,
			setupMocks: func(r *mocks.ProductRepository) {
				r.On("AddStock", mock.Anything, int64(1), int64(10)).Return(errors.New("failed to add stock"))
			},
			expectedError: errors.New("failed to add stock"),
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}

			tc.setupMocks(repository)

			productService := service.NewProductService(repository, nil, nil, "")

			err := productService.AddStock(context.Background(), tc.productID, tc.quantity)

			assert.Equal(t, tc.expectedError, err)

			repository.AssertExpectations(t)
		})
	}
}

func TestReduceStock(t *testing.T) {
	testCases := []struct {
		name          string
		productID     int64
		quantity      int64
		setupMocks    func(r *mocks.ProductRepository)
		expectedError error
	}{
		{
			name:      "Reduce stock success",
			productID: 1,
			quantity:  10,
			setupMocks: func(r *mocks.ProductRepository) {
				r.On("ReduceStock", mock.Anything, int64(1), int64(10)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Failed to reduce stock",
			productID: 1,
			quantity:  10,
			setupMocks: func(r *mocks.ProductRepository) {
				r.On("ReduceStock", mock.Anything, int64(1), int64(10)).Return(errors.New("failed to reduce stock"))
			},
			expectedError: errors.New("failed to reduce stock"),
		},
	}

	logger.Init("info")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := &mocks.ProductRepository{}

			tc.setupMocks(repository)

			productService := service.NewProductService(repository, nil, nil, "")

			err := productService.ReduceStock(context.Background(), tc.productID, tc.quantity)

			assert.Equal(t, tc.expectedError, err)

			repository.AssertExpectations(t)
		})
	}
}
