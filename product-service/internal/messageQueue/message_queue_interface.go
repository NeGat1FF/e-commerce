package messagequeue

import (
	"context"
)

// MessageQueue defines the methods that any
// message queue provider needs to implement to publish messages.
type MessageQueue interface {
	// PublishMessage publishes a message to a queue.
	PublishMessage(ctx context.Context, exchangeName, routingKey string, message interface{}) error
}
