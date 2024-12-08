package messagequeue

import (
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

func (r *RabbitMQClient) Close() error {
	return r.ch.Close()
}

func (r *RabbitMQClient) GetConsumeChannel(queueName string, consumer string, autoAck bool, exclusive bool, noWait bool) (<-chan amqp.Delivery, error) {
	return r.ch.Consume(queueName, consumer, autoAck, exclusive, false, noWait, nil)
}
