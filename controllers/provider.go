package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"io"
	"time"

	"net/http"

	"provider_gateway_mq/consts"
	"provider_gateway_mq/schema"
	"provider_gateway_mq/utils"
)

func Provider(c *gin.Context) {
	requestIdHeaderName := consts.RequestIdHttpHeaderName
	requestId := c.GetString(requestIdHeaderName)

	var msg schema.BodyRequest

	if binderr := c.ShouldBindJSON(&msg); binderr != nil {
		log.Error().Err(binderr).Str("request_id", requestId).
			Msg("Error occured while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	headers := make(amqp.Table)
	for k, v := range c.Request.Header {
		headers[k] = v[0]
	}

	msgCreate := &schema.CreateMessage{
		Body:          msg,
		CorrelationId: utils.GetCorrelationId(),
		Headers:       headers,
		ReplyTo:       consts.AnswerQueueName,
	}

	// Create channel and add rchans with CorrId
	replyChannel := make(chan schema.ReplyMessage)
	ReplyChannels[msgCreate.CorrelationId] = replyChannel

	PublishChannels <- *msgCreate

	waitReply(msgCreate.CorrelationId, replyChannel, c.Writer)

}

func waitReply(CorrelationId string, replyChannel chan schema.ReplyMessage, w gin.ResponseWriter) {
	for {
		select {
		case msgReply := <-replyChannel:
			// response received
			log.Printf("INFO: received reply: %v correlationId: %s", msgReply.Data, CorrelationId)
			for k := range msgReply.Headers {
				if str, ok := msgReply.Headers[k].(string); ok {
					w.Header().Set(k, str)
				}
			}
			log.Printf("INFO: received reply: %v requestId: %s",
				msgReply.Data, msgReply.Headers[consts.RequestIdHttpHeaderName])
			// send response back to client

			response(w, msgReply.Data, 200)
			delete(ReplyChannels, CorrelationId)

			fmt.Printf("LEN RCHANS afrer delete %v", len(ReplyChannels))
			return

		case <-time.After(90 * time.Second):
			// timeout
			log.Printf("ERROR: request timeout msg with CorrelationId %s", CorrelationId)

			response(w, "Timeout", 408)

			delete(ReplyChannels, CorrelationId)
			return
		}
	}

}

func response(w gin.ResponseWriter, result string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, result)
}
