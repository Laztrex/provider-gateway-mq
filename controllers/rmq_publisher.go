package controllers

import (
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"provider_gateway_mq/utils"
)

type RMQProducer struct {
	Queue            string
	ConnectionString string
}

func (x RMQProducer) OnError(err error, msg string) {
	if err != nil {
		log.Err(err).Msgf("Error while publishing message on '$s' queue. Error message: %s", x.Queue, msg)
	}
}

func (x RMQProducer) PublishMessages() {
	tlsConf := utils.GetRmqTlsConf()

	exchange := utils.GetEnvVar("TOPIC")
	routingKey := utils.GetEnvVar("ROUTING_KEY")

	conn, err := amqp.DialTLS(x.ConnectionString, tlsConf)
	x.OnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	channel, err := conn.Channel()
	x.OnError(err, "Failed to open a channel")

	defer channel.Close()

	err = channel.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	queue, err := channel.QueueDeclare(
		x.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	x.OnError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		queue.Name,
		utils.GetEnvVar("BINDING_KEY"),
		exchange,
		false,
		nil,
	)
	x.OnError(err, "Failed to bind a queue")

	for {
		select {
		case msg := <-PublishChannels:
			err = channel.Publish(
				exchange,
				routingKey,
				false,
				false,
				amqp.Publishing{
					ContentType:   "application/json",
					Body:          []byte(msg.Body.Message),
					Headers:       msg.Headers,
					CorrelationId: msg.CorrelationId,
					ReplyTo:       msg.ReplyTo,
				},
			)
			x.OnError(err, "Failed to publish a message")

			log.Printf("INFO: [%v] - published msg: %v", msg.CorrelationId, msg.Body)
		}
	}
}
