package schema

import "github.com/streadway/amqp"

type BodyRequest struct {
	Message string `json:"data" binding:"required"`
}

type CreateMessage struct {
	Body          BodyRequest `json:"data"`
	CorrelationId string
	Headers       amqp.Table
	ReplyTo       string
}

type ReplyMessage struct {
	Data          string `json:"data" binding:"required"`
	CorrelationId string
	Headers       amqp.Table
}
