package controllers

import (
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

func (conn *RMQSpec) PublishDeclare() {
	var err error

	if conn.Exchange != "" {
		err = conn.ExchangeDeclare()
		conn.OnError(err, "Failed to declare exchange while publishing")
	}

	err = conn.QueueDeclare()
	conn.OnError(err, "Failed to declare a queue while publishing")

	if conn.Exchange != "" {
		err = conn.QueueBind()
		conn.OnError(err, "Failed to bind a queue while publishing")
	}
}

func (conn *RMQSpec) PublishMessages() {

	for {
		select {
		case err := <-conn.Err:
			err = conn.Reconnect()
			if err != nil {
				panic(err)
			}

		case msg := <-PublishChannels:
			err := conn.Channel.Publish(
				conn.Exchange, // exchange
				//conn.RoutingKey, // routing key
				msg.RoutingKey,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					Body:          []byte(msg.Body.Message),
					Headers:       msg.Headers,
					CorrelationId: msg.CorrelationId,
					ReplyTo:       conn.ReplyTo,
				},
			)
			if err != nil {
				log.Err(err).Msgf("ERROR: fail to publish msg: %s", msg.CorrelationId)
			}
			log.Printf("INFO: [%v] - published", msg.CorrelationId)
		}
	}
}
