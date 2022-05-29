package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"provider_gateway_mq/consts"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestId string

		rqUID := c.Request.Header.Get(consts.RequestIdHttpHeaderName)
		c.Request.Header.Del(consts.RequestIdHttpHeaderName)

		if rqUID != "" {
			requestId = rqUID
		} else {
			requestId = uuid.New().String()
		}

		c.Set("RqUID", requestId)
		c.Request.Header["RqUID"] = []string{requestId}
		c.Next()
	}
}
