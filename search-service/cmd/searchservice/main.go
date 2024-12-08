package main

import (
	"github.com/NeGat1FF/e-commerce/search-service/internal/config"
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

	cfg := config.LoadConfig()

	esCfg := es.Config{
		Addresses: []string{
			cfg.ElasticURL,
		},
		Username: cfg.ElasticUsername,
		Password: cfg.ElasticPassword,
	}

	// Connect to Elasticsearch
	esClient, err := elasticsearch.NewElasticClient(esCfg, "products")
	if err != nil {
		panic(err)
	}

	// Connect to RabbitMQ
	conn, err := messagequeue.ConnectRabbitMQ(cfg.MessageBrokerURL)
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
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
