package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"io"
	"net/http"
	"time"

	"gateway_mq/internal/consts"
	"gateway_mq/internal/controllers"
	"gateway_mq/internal/schemas"
	"gateway_mq/internal/utils"
)

func Provider(c *gin.Context) {
	var msg schemas.MessageRequest
	//var requestId string

	requestIdHeaderName := consts.KeyRequestId
	requestId := c.GetString(requestIdHeaderName)
	log.Info().Msgf("requestId in Provider: %s", requestId)

	if binderr := c.ShouldBindJSON(&msg.Message); binderr != nil {
		log.Error().Err(binderr).Str(requestIdHeaderName, requestId).
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
		//CorrelationId: utils.GetCorrelationId(),
		CorrelationId: utils.SetCorrelationId(requestId),
		Body:          msg,
		Headers:       headers,
		RoutingKey:    c.Request.Header.Get(consts.KeyRoutingKey),
	}

	replyChannel := make(chan schemas.MessageReply)
	controllers.ReplyChannels[msgCreate.CorrelationId] = replyChannel

	controllers.PublishChannels <- *msgCreate

	waitReply(c, msgCreate.CorrelationId, replyChannel, c.Writer) //, w http.ResponseWriter

}

func waitReply(c *gin.Context, correlationId string, replyChannel chan schemas.MessageReply, w gin.ResponseWriter) {
	for {
		select {
		case msgReply := <-replyChannel:

			for k := range msgReply.Headers {
				if str, ok := msgReply.Headers[k].(string); ok {
					w.Header().Set(k, str)
				}
			}

			//rcvRequestId := w.Header().Get(consts.KeyRequestId)
			//c.Set(consts.KeyRequestId, rcvRequestId)

			response(w, msgReply.Data, 200)

			delete(controllers.ReplyChannels, correlationId)
			return

		case <-time.After(90 * time.Second):

			log.Printf("ERROR: request timeout msg with correlation id: %s", correlationId)

			response(w, "Timeout", 408)

			delete(controllers.ReplyChannels, correlationId)
			return
		}
	}
}

func response(w gin.ResponseWriter, resp string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, resp)
}
