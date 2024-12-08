package main

import (
	"os"

	"github.com/NeGat1FF/e-commerce/search-service/internal/elasticsearch"
	"github.com/NeGat1FF/e-commerce/search-service/internal/handlers"
	messagequeue "github.com/NeGat1FF/e-commerce/search-service/internal/messageQueue"
	"github.com/NeGat1FF/e-commerce/search-service/internal/service"
	"github.com/NeGat1FF/e-commerce/search-service/logger"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger.Init("info")

	cfg := es.Config{
		APIKey:   os.Getenv("ELASTICSEARCH_URL"),
		Username: "search_service",
		Password: "search_service",
	}

	// Connect to Elasticsearch
	esClient, err := elasticsearch.NewElasticClient(cfg, "products")
	if err != nil {
		panic(err)
	}

	// Connect to RabbitMQ
	conn, err := messagequeue.ConnectRabbitMQ(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err)
	}

	rmqClient, err := messagequeue.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	// Create a new SearchService
	searchService := service.NewSearchService(rmqClient, esClient, "product")

	// Start the SearchService
	err = searchService.Start()
	if err != nil {
		panic(err)
	}

	// Create a new HTTP server
	router := gin.Default()

	// Create a new SearchHandler
	searchHandler := handlers.NewSearchHandler(searchService)

	// Register the SearchHandler
	router.GET("/api/v1/products/search", searchHandler.SearchProducts)

	// Start the HTTP server
	err = router.Run(":8090")
	if err != nil {
		panic(err)
	}
}
