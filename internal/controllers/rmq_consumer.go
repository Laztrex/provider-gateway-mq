package controllers

import (
	"github.com/rs/zerolog/log"

	"gateway_mq/internal/consts"
	"gateway_mq/internal/schemas"
)

func (conn *RMQSpec) ConsumeDeclare() {

	err := conn.QueueDeclare()
	conn.OnError(err, "Failed to declare a queue while consuming")

}

func (conn *RMQSpec) ConsumeMessages() {

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
		case err := <-conn.Err:
			err = conn.Reconnect()
			if err != nil {
				panic(err)
			}

		case msg := <-msgChannel:

			log.Info().
				Str(consts.KeyCorrelationId, msg.CorrelationId).
				Msgf("CONSUME: %s", conn.Queue)

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
				log.Error().Err(err).
					Str(consts.KeyCorrelationId, msg.CorrelationId).
					Msgf("ERROR: fail to ack: %s", err.Error())
			}

			if rchan, ok := ReplyChannels[msgRply.CorrelationId]; ok {
				rchan <- *msgRply
			}
		}
	}
}
