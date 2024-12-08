package messagequeue

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func ConnectRabbitMQ(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewRabbitMQClient(conn *amqp.Connection) (*RabbitMQClient, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{
		conn: conn,
		ch:   channel,
	}, nil
}

func (r *RabbitMQClient) PublishMessage(ctx context.Context, exchangeName, routingKey string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.ch.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}
