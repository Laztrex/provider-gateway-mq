package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	stdlog "log"
	"os"
	"time"

	"gateway_mq/internal/consts"
)

// RequestLogger Request logging middleware
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestId := c.Request.Header.Get(consts.KeyRequestId)
		c.Request.Header.Del(consts.KeyRequestId)

		if requestId == "" {
			requestId = uuid.New().String()
		}

		c.Set(consts.KeyRequestId, requestId)
		c.Request.Header[consts.KeyRequestId] = []string{requestId}

		method := c.Request.Method
		path := c.Request.URL.Path
		t := time.Now()

		c.Next()

		latency := float32(time.Since(t).Seconds())
		status := c.Writer.Status()

		stdlog.Printf("%v Request: %s - %s - %f - [%s]", status, method, path, latency, requestId)
		logToFile(status, requestId, method, path, latency)
	}
}

func logToFile(status int, requestId string, method string, path string, latency float32) {
	tempFile, err := os.OpenFile(consts.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("there was an error creating a temporary file four our log")
	}

	fileLogger := zerolog.New(tempFile)
	fileLogger.Info().
		Int("status", status).
		Str(consts.KeyRequestId, requestId).
		Str("method", method).
		Str("path", path).
		Float32("latency", latency).
		Msg("SUCCESS")
}
