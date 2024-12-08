package messagequeue_test

import (
	"context"
	"testing"

	messagequeue "github.com/NeGat1FF/e-commerce/search-service/internal/messageQueue"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"

	mq "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func setupTestMessageQueue(t *testing.T, ctx context.Context) (*amqp.Connection, func()) {

	rabbitmqContainer, err := mq.Run(ctx, "rabbitmq:3.7.25-management-alpine", mq.WithAdminUsername("admin"), mq.WithAdminPassword("admin"))
	require.NoError(t, err)

	url, err := rabbitmqContainer.AmqpURL(ctx)
	require.NoError(t, err)

	// Set up RabbitMQ connection
	conn, err := amqp.Dial(url)
	require.NoError(t, err)

	// Set up RabbitMQ channel
	ch, err := conn.Channel()
	require.NoError(t, err)

	// Declare the exchange
	err = ch.ExchangeDeclare("testExchange", "topic", true, false, false, false, nil)
	require.NoError(t, err)

	// Declare the queue
	_, err = ch.QueueDeclare("test_queue", true, false, false, false, nil)
	require.NoError(t, err)

	// Bind the queue to the exchange
	err = ch.QueueBind("test_queue", "product.test", "testExchange", false, nil)
	require.NoError(t, err)

	cleanup := func() {
		rabbitmqContainer.Terminate(ctx)
		conn.Close()
	}

	return conn, cleanup
}

func TestRabbitMQClient_ConsumeMessage(t *testing.T) {
	// Set up RabbitMQ connection and channel
	conn, cleanup := setupTestMessageQueue(t, context.Background())
	defer cleanup()

	// Create a new RabbitMQ client
	rmqClient, err := messagequeue.NewRabbitMQClient(conn)
	require.NoError(t, err)

	ch, err := conn.Channel()
	require.NoError(t, err)

	// Publish the message
	err = ch.Publish("testExchange", "product.test", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(`{"key":"value"}`),
	})
	require.NoError(t, err)

	// Create a new channel
	msgs, err := rmqClient.GetConsumeChannel("test_queue", "test", false, false, false)
	require.NoError(t, err)

	// Check if the message is the same
	msg := <-msgs
	require.Equal(t, `{"key":"value"}`, string(msg.Body))
}
