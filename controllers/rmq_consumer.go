package controllers

import (
	"github.com/rs/zerolog/log"

	"provider_gateway_mq/schemas"
)

func (conn *RMQSpec) ConsumeMessages() {
	// conn

	err := conn.Connect()
	conn.OnError(err, "Failed to connect to RabbitMQ while consuming")

	err = conn.QueueDeclare()
	conn.OnError(err, "Failed to declare a queue while consuming")

	msgChannel, err := conn.Channel.Consume(
		conn.Queue, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	conn.OnError(err, "ERROR: fail create channel")

	for {
		select {
		case msg := <-msgChannel:

			if msg.CorrelationId == "" {
				continue
			}

			msgRply := &schemas.MessageReply{
				CorrelationId: msg.CorrelationId,
				Data:          string(msg.Body),
				Headers:       msg.Headers,
			}

			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: fail to ack: %s", err.Error())
			}

			if rchan, ok := ReplyChannels[msgRply.CorrelationId]; ok {
				rchan <- *msgRply
			}
		}
	}
}
