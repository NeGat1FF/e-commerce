package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/NeGat1FF/e-commerce/product-service/internal/cache"
	messagequeue "github.com/NeGat1FF/e-commerce/product-service/internal/messageQueue"
	"github.com/NeGat1FF/e-commerce/product-service/internal/models"
	"github.com/NeGat1FF/e-commerce/product-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/product-service/pkg/logger"
	"github.com/NeGat1FF/e-commerce/product-service/proto"
	"go.uber.org/zap"
)

var (
	ErrFailedToCreateProduct  = errors.New("failed to create product")
	ErrFailedToUpdateProduct  = errors.New("failed to update product")
	ErrFailedToDeleteProduct  = errors.New("failed to delete product")
	ErrFailedToAddStock       = errors.New("failed to add stock")
	ErrFailedToReduceStock    = errors.New("failed to reduce stock")
	ErrProductNotFound        = errors.New("product not found")
	ErrFailedToGetProductByID = errors.New("failed to get product by id")
	ErrProductAlreadyExists   = errors.New("product already exists")
)

type ProductService struct {
	proto.UnimplementedPriceServiceServer
	repo         repository.ProductRepository
	messageQueue messagequeue.MessageQueue
	cache        cache.Cache
	exchangeName string
}

func NewProductService(repo repository.ProductRepository, messageQueue messagequeue.MessageQueue, cache cache.Cache, exchangeName string) *ProductService {
	return &ProductService{
		repo:         repo,
		messageQueue: messageQueue,
		cache:        cache,
		exchangeName: exchangeName,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, product models.Product) error {
	logger.Logger.Info("Creating product", zap.Any("product", product))
	err := ps.repo.CreateProduct(ctx, product)
	if err != nil {
		logger.Logger.Error("Failed to create product", zap.Error(err))
		if err == repository.ErrProductAlreadyExists {
			return ErrProductAlreadyExists
		}
		return ErrFailedToCreateProduct
	}
	logger.Logger.Info("Product created successfully")

	go func() {
		err := ps.messageQueue.PublishMessage(ctx, ps.exchangeName, "product.created", product)
		if err != nil {
			logger.Logger.Error("Failed to publish message", zap.Error(err))
		}
	}()

	return nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, id int64, updateFields map[string]interface{}) error {
	logger.Logger.Info("Updating product", zap.Int64("id", id), zap.Any("updateFields", updateFields))
	err := ps.repo.UpdateProduct(ctx, id, updateFields)
	if err != nil {
		logger.Logger.Error("Failed to update product", zap.Error(err))
		if err == repository.ErrProductNotFound {
			return ErrProductNotFound
		}
		return ErrFailedToUpdateProduct
	}
	logger.Logger.Info("Product updated successfully")

	go func() {
		// Invalidate cache
		err := ps.cache.Del(ctx, fmt.Sprintf("products:%d", id))
		if err != nil {
			logger.Logger.Error("Failed to delete product from cache", zap.Error(err))
		}

		// Add id to updateFields
		updateFields["id"] = id

		err = ps.messageQueue.PublishMessage(ctx, ps.exchangeName, "product.updated", updateFields)
		if err != nil {
			logger.Logger.Error("Failed to publish message", zap.Error(err))
		}
	}()

	return nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	logger.Logger.Info("Deleting product", zap.Int64("id", id))
	err := ps.repo.DeleteProduct(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to delete product", zap.Error(err))
		if err == repository.ErrProductNotFound {
			return ErrProductNotFound
		}
		return ErrFailedToDeleteProduct
	}
	logger.Logger.Info("Product deleted successfully")

	go func() {
		// Delete product from cache
		err := ps.cache.Del(ctx, fmt.Sprintf("products:%d", id))
		if err != nil {
			logger.Logger.Error("Failed to delete product from cache", zap.Error(err))
		}

		err = ps.messageQueue.PublishMessage(ctx, ps.exchangeName, "product.deleted", map[string]int64{"id": id})
		if err != nil {
			logger.Logger.Error("Failed to publish message", zap.Error(err))
		}
	}()

	return nil
}

func (ps *ProductService) GetProductByID(ctx context.Context, id int64) (models.UserProduct, error) {
	var product models.UserProduct
	err := ps.cache.Get(ctx, fmt.Sprintf("products:%d", id), &product)
	if err == nil {
		return product, nil
	}

	product, err = ps.repo.GetProductByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to get product by id", zap.Error(err))
		if err == repository.ErrProductNotFound {
			return product, ErrProductNotFound
		}
		return product, ErrFailedToGetProductByID
	}

	go func() {
		err := ps.cache.Set(ctx, fmt.Sprintf("products:%d", id), product)
		if err != nil {
			logger.Logger.Error("Failed to add product to cache", zap.Error(err))
		}
	}()

	return product, nil
}

func (ps *ProductService) GetProductsByCategory(ctx context.Context, category string, page, limit int) ([]models.UserProduct, error) {
	return ps.repo.GetProductsByCategory(ctx, category, page, limit)
}

func (ps *ProductService) GetStock(ctx context.Context, id int64) (int64, error) {
	return ps.repo.GetStock(ctx, id)
}

func (ps *ProductService) AddStock(ctx context.Context, id int64, quantity int64) error {
	logger.Logger.Info("Adding stock", zap.Int64("id", id), zap.Int64("quantity", quantity))
	err := ps.repo.AddStock(ctx, id, quantity)
	if err != nil {
		logger.Logger.Error("Failed to add stock", zap.Error(err))
		return ErrFailedToAddStock
	}
	logger.Logger.Info("Stock added successfully")

	return nil
}

func (ps *ProductService) ReduceStock(ctx context.Context, id int64, quantity int64) error {
	logger.Logger.Info("Reducing stock", zap.Int64("id", id), zap.Int64("quantity", quantity))
	err := ps.repo.ReduceStock(ctx, id, quantity)
	if err != nil {
		logger.Logger.Error("Failed to reduce stock", zap.Error(err))
		return err
	}
	logger.Logger.Info("Stock reduced successfully")

	return nil
}

func (ps *ProductService) GetPrice(ctx context.Context, in *proto.PriceRequest) (*proto.PriceResponse, error) {
	productID, err := strconv.ParseInt(in.ProductId, 10, 64)
	if err != nil {
		logger.Logger.Error("Failed to parse product id", zap.Error(err))
		return nil, err
	}

	product, err := ps.GetProductByID(ctx, productID)
	if err != nil {
		logger.Logger.Error("Failed to get product by id", zap.Error(err))
		return nil, err
	}

	return &proto.PriceResponse{
		Price: fmt.Sprintf("%.2f", product.Price),
	}, nil
}
