package schemas

import (
	"github.com/streadway/amqp"
)

// MessageRequest Message is RequestBody from AC
type MessageRequest struct {
	Message map[string]interface{}
}

type MessageCreate struct {
	CorrelationId string
	Body          MessageRequest
	Headers       amqp.Table
	RoutingKey    string
}

// MessageReply is Response mapping queue
type MessageReply struct {
	CorrelationId string
	Data          string `json:"data" binding:"required"`
	Headers       amqp.Table
}
