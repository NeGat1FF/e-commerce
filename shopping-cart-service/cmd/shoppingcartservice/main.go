package main

import (
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/config"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/db"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/handlers"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/service"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config.InitConfig()

	db, err := db.InitDB(config.GetConfig().DB_URL)
	if err != nil {
		panic(err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcConn, err := grpc.NewClient(config.GetConfig().PRICE_SERVICE, opts...)
	if err != nil {
		panic(err)
	}

	repo := repository.NewShoppingCartRepo(db)
	priceService := proto.NewPriceServiceClient(grpcConn)

	service := service.NewCartService(repo, priceService)

	handler := handlers.NewCartHandler(service)

	ginServer := gin.Default()

	group := ginServer.Group("/api/v1/cart")

	group.Use(handlers.AuthMiddleware())

	group.POST("", handler.AddItem)
	group.PUT("/:itemID", handler.SetQuantity)
	group.DELETE("/:itemID", handler.RemoveItem)
	group.GET("", handler.GetItems)
	group.DELETE("", handler.DeleteCart)

	err = ginServer.Run(":" + config.GetConfig().SERVER_PORT)
	if err != nil {
		panic(err)
	}
}
