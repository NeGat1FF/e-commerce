package main

import (
	"net/http"

	"github.com/NeGat1FF/e-commerce/order-service/internal/config"
	"github.com/NeGat1FF/e-commerce/order-service/internal/db"
	"github.com/NeGat1FF/e-commerce/order-service/internal/handlers"
	"github.com/NeGat1FF/e-commerce/order-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/order-service/internal/service"
	"github.com/NeGat1FF/e-commerce/order-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	config.InitConfig()

	db, err := db.InitDB(config.GetConfig().DB_URL)
	if err != nil {
		panic(err)
	}

	repo := repository.NewOrderRepo(db)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcClient, err := grpc.NewClient(config.GetConfig().PRICE_SERVICE, opts...)
	if err != nil {
		panic(err)
	}

	priceService := proto.NewPriceServiceClient(grpcClient)

	orderService := service.NewOrderService(repo, priceService, config.GetConfig().JWTSecret)

	// Create a new OrderHandler
	orderHandler := handlers.NewOrderHandler(orderService)

	// Create a new server mux
	mux := http.NewServeMux()

	// Register the handler functions
	mux.HandleFunc("POST api/v1/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET api/v1/orders", orderHandler.GetOrders)
	mux.HandleFunc("GET api/v1/orders/{id}", orderHandler.GetOrder)
	mux.HandleFunc("PUT api/v1/orders/{id}", orderHandler.UpdateOrder)

	// Create a new server
	err = http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		mux.ServeHTTP(w, r)
	}))

	if err != nil {
		panic(err)
	}
}
