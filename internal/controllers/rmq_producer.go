package controllers

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	"gateway_mq/internal/consts"
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
			log.Info().
				Str(consts.KeyCorrelationId, msg.CorrelationId).
				Msgf("PRODUCE: %s", conn.Queue)

			body, err := json.Marshal(msg.Body.Message)
			err = conn.Channel.Publish(
				conn.Exchange, // exchange
				//conn.RoutingKey, // routing key
				msg.RoutingKey,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					Body:          body,
					Headers:       msg.Headers,
					CorrelationId: msg.CorrelationId,
					ReplyTo:       conn.ReplyTo,
				},
			)
			if err != nil {
				log.Error().Err(err).
					Str(consts.KeyCorrelationId, msg.CorrelationId).
					Msgf("ERROR: fail to publish msg: %s", msg.CorrelationId)
			}
			log.Info().
				Str(consts.KeyCorrelationId, msg.CorrelationId).
				Msgf("INFO: [%v] - published", msg.CorrelationId)
		}
	}
}
