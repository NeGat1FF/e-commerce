package main

import (
	"fmt"
	"net"

	"github.com/NeGat1FF/e-commerce/product-service/internal/cache"
	"github.com/NeGat1FF/e-commerce/product-service/internal/config"
	"github.com/NeGat1FF/e-commerce/product-service/internal/db"
	"github.com/NeGat1FF/e-commerce/product-service/internal/handlers"
	messagequeue "github.com/NeGat1FF/e-commerce/product-service/internal/messageQueue"
	"github.com/NeGat1FF/e-commerce/product-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/product-service/internal/service"
	"github.com/NeGat1FF/e-commerce/product-service/pkg/logger"
	"github.com/NeGat1FF/e-commerce/product-service/proto"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	config := config.LoadConfig()

	logger.Init(config.LogLevel)

	db, err := db.InitDB(config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	repo := repository.NewMongoRepository(db.Database("product").Collection("products"))

	opts, err := redis.ParseURL(config.CacheURL)
	if err != nil {
		panic(err)
	}
	cache := cache.NewRedisClient(redis.NewClient(opts))

	conn, err := messagequeue.ConnectRabbitMQ(config.MessageBrokerURL)
	if err != nil {
		panic(err)
	}
	mqClient, err := messagequeue.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	// Initialize the service
	service := service.NewProductService(repo, mqClient, cache, config.MessageBrokerExchange)

	s := grpc.NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCPort))
	if err != nil {
		logger.Logger.Error("Failed to listen", zap.Error(err))
		panic(err)
	}

	proto.RegisterPriceServiceServer(s, service)

	go func() {
		logger.Logger.Info("Starting gRPC server")
		err := s.Serve(lis)
		if err != nil {
			logger.Logger.Error("Failed to start gRPC server", zap.Error(err))
		}
	}()

	// Initialize the handlers
	productHandler := handlers.NewProductHandler(service)

	ginServer := gin.New()
	ginServer.Use(handlers.Logging())
	ginServer.Use(handlers.ErrorHandling())
	ginServer.Use(gin.Recovery())

	group := ginServer.Group("/api/v1/products")
	group.POST("/", handlers.Auth(), handlers.ValidateProduct(), productHandler.CreateProduct)
	group.PUT("/:id", handlers.Auth(), productHandler.UpdateProduct)
	group.DELETE("/:id", handlers.Auth(), productHandler.DeleteProduct)

	group.POST("/:id/add-stock", handlers.Auth(), productHandler.AddStock)
	group.POST("/:id/reduce-stock", handlers.Auth(), productHandler.ReduceStock)
	group.GET("/:id/stock", productHandler.GetStock)

	group.GET("/:id", productHandler.GetProductByID)
	group.GET("/", productHandler.GetProductsByCategory)

	ginServer.Run(fmt.Sprintf(":%s", config.ServerPort))
}
