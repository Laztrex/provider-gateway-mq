package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"provider_gateway_mq/consts"
)

// RequestID Request ID middleware
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestId string

		requestId = c.Request.Header.Get(consts.RequestIdHttpHeaderName)
		c.Request.Header.Del(consts.RequestIdHttpHeaderName)

		if requestId == "" {
			requestId = uuid.New().String()
		}

		// Set context variable
		c.Set("requestId", requestId)
		c.Request.Header[consts.RequestIdHttpHeaderName] = []string{requestId}

		c.Next()
	}
}
