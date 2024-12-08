package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/NeGat1FF/e-commerce/search-service/internal/elasticsearch"
	messagequeue "github.com/NeGat1FF/e-commerce/search-service/internal/messageQueue"
	"github.com/NeGat1FF/e-commerce/search-service/internal/models"
	"github.com/NeGat1FF/e-commerce/search-service/logger"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func handleError(msg amqp091.Delivery, err error, errMsg string) {
	if !msg.Redelivered {
		logger.Logger.Error(errMsg, zap.Error(err))
		msg.Nack(false, true)
	} else {
		logger.Logger.Error(errMsg+" twice, rejecting...", zap.Error(err))
		msg.Reject(false)
	}
}

type SearchService struct {
	messageQueue *messagequeue.RabbitMQClient
	elastic      *elasticsearch.ElasticClient
	queueName    string
}

func NewSearchService(mq *messagequeue.RabbitMQClient, elastic *elasticsearch.ElasticClient, queueName string) *SearchService {
	return &SearchService{
		messageQueue: mq,
		elastic:      elastic,
		queueName:    queueName,
	}
}

func (s *SearchService) Start() error {
	ch, err := s.messageQueue.GetConsumeChannel(s.queueName, "", false, false, false)
	if err != nil {
		return err
	}
	go s.processMessages(ch)
	return nil
}

func (s *SearchService) Close() {
	s.messageQueue.Close()
}

func (s *SearchService) SearchProducts(ctx context.Context, query map[string]string, sortOrder string, min_price, max_price, page, limit int) ([]models.Product, error) {
	models, err := s.elastic.SearchProducts(ctx, query, sortOrder, min_price, max_price, page, limit)
	if err != nil {
		logger.Logger.Error("Failed to search products", zap.Error(err))
		return nil, err
	}
	return models, nil
}

func (s *SearchService) indexProduct(msg amqp091.Delivery) {
	var product models.Product
	// Read the message body
	err := json.Unmarshal(msg.Body, &product)
	if err != nil {
		handleError(msg, err, "Failed to unmarshal message")
		return
	}

	// Index the product
	err = s.elastic.IndexProduct(context.Background(), product)
	if err != nil {
		handleError(msg, err, "Failed to index product")
		return
	}

	logger.Logger.Info("Product indexed successfully", zap.Int64("id", product.ID))
	msg.Ack(false) // Acknowledge the message
}

func (s *SearchService) updateProduct(msg amqp091.Delivery) {
	// Read the message body
	var updateFields map[string]interface{}
	err := json.Unmarshal(msg.Body, &updateFields)
	if err != nil {
		handleError(msg, err, "Failed to unmarshal message")
		return
	}

	// Get the product ID
	idf, ok := updateFields["id"].(float64)
	if !ok {
		logger.Logger.Error("Failed to get product ID from message")
		msg.Reject(false)
		return
	}

	id := int64(idf)

	// Update the product
	err = s.elastic.UpdateProduct(context.Background(), id, updateFields)
	if err != nil {
		handleError(msg, err, "Failed to update product")
		return
	}

	logger.Logger.Info("Product updated successfully", zap.Int64("id", id))
	msg.Ack(false) // Acknowledge the message
}

func (s *SearchService) deleteProduct(msg amqp091.Delivery) {
	// Get the product ID from the message body
	var idMap map[string]int64
	err := json.Unmarshal(msg.Body, &idMap)
	if err != nil {
		handleError(msg, err, "Failed to unmarshal message")
		return
	}

	id := idMap["id"]

	// Delete the product
	err = s.elastic.DeleteProduct(context.Background(), id)
	if err != nil {
		handleError(msg, err, "Failed to delete product")
		return
	}

	logger.Logger.Info("Product deleted successfully", zap.Int64("id", id))
	msg.Ack(false) // Acknowledge the message
}

func (s *SearchService) processMessages(ch <-chan amqp091.Delivery) {
	for msg := range ch {
		// Get the route from the message body
		route := strings.Split(string(msg.RoutingKey), ".")[1]

		switch route {
		case "created":
			s.indexProduct(msg)
		case "updated":
			s.updateProduct(msg)
		case "deleted":
			s.deleteProduct(msg)
		default:
			logger.Logger.Error("Invalid route", zap.String("route", route))
			msg.Reject(false)
		}
	}
}
