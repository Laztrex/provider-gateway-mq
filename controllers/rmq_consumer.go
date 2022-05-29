package controllers

import (
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"os"
	"provider_gateway_mq/schema"
	"provider_gateway_mq/utils"
)

type RMQConsumer struct {
	Queue            string
	ConnectionString string
}

func (x RMQConsumer) OnError(err error, msg string) {
	if err != nil {
		log.Err(err).Msgf("Error while consuming message on '$s' queue. Error message: %s", x.Queue, msg)
		os.Exit(1)
	}
}

func (x RMQConsumer) ConsumeMessages() {
	tlsConf := utils.GetRmqTlsConf()

	// set conn
	conn, err := amqp.DialTLS(x.ConnectionString, tlsConf)
	x.OnError(err, "Failed to connect to RabbitMQ")
	// notify ?

	defer conn.Close()

	log.Printf("INFO: Successful init consumer conn")

	amqpChannel, err := conn.Channel()
	x.OnError(err, "Failed to open a channel")

	queue, err := amqpChannel.QueueDeclare(
		x.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	x.OnError(err, "Failed to declare a queue")

	msgChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	x.OnError(err, "ERROR: Failed to consume channel")

	for {
		select {
		case msg := <-msgChannel:
			if msg.CorrelationId == "" {
				continue
			}

			replyMsg := &schema.ReplyMessage{
				CorrelationId: msg.CorrelationId,
				Data:          string(msg.Body),
				Headers:       msg.Headers,
			}
			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: Failed to ack message", err.Error())
			}

			// find waiting channel(with corr_id) and forward the reply to it
			if replyChan, ok := ReplyChannels[replyMsg.CorrelationId]; ok {
				replyChan <- *replyMsg
			}

		}
	}
}
