package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"io"
	"net/http"
	"time"

	"provider_gateway_mq/consts"
	"provider_gateway_mq/schemas"
	"provider_gateway_mq/utils"
)

func Provider(c *gin.Context) {
	var msg schemas.MessageRequest

	requestIdHeaderName := consts.RequestIdHttpHeaderName
	requestId := c.GetString(requestIdHeaderName)

	if binderr := c.ShouldBindJSON(&msg); binderr != nil {
		log.Error().Err(binderr).Str("request_id", requestId).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	headers := make(amqp.Table)

	for k, v := range c.Request.Header {
		headers[k] = v[0]
		//headers[textproto.CanonicalMIMEHeaderKey(k)] = v[0]
	}

	msgCreate := &schemas.MessageCreate{
		CorrelationId: utils.GetCorrelationId(),
		Body:          msg,
		Headers:       headers,
		ReplyTo:       consts.AnswerQueueName,
		RoutingKey:    c.Request.Header.Get("routing-key"),
	}

	replyChannel := make(chan schemas.MessageReply)
	ReplyChannels[msgCreate.CorrelationId] = replyChannel

	PublishChannels <- *msgCreate

	waitReply(msgCreate.CorrelationId, replyChannel, c.Writer) //, w http.ResponseWriter

}

func waitReply(CorrelationId string, replyChannel chan schemas.MessageReply, w gin.ResponseWriter) {
	for {
		select {
		case msgReply := <-replyChannel:

			log.Printf("INFO: [%s] received", CorrelationId)

			for k := range msgReply.Headers {
				if str, ok := msgReply.Headers[k].(string); ok {
					w.Header().Set(k, str)
				}
			}

			response(w, msgReply.Data, 200)

			delete(ReplyChannels, CorrelationId)
			return

		case <-time.After(90 * time.Second):

			log.Printf("ERROR: request timeout msg with correlation id: %s", CorrelationId)

			response(w, "Timeout", 408)

			delete(ReplyChannels, CorrelationId)
			return
		}
	}
}

func response(w gin.ResponseWriter, resp string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, resp)
}
