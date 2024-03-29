package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gateway_mq/internal/middlewares"
	"gateway_mq/internal/transport"
)

// SetupApp Function to setup the app object
func SetupApp() *gin.Engine {
	log.Info().Msg("Initializing service")

	app := gin.New()

	app.Use(gin.Recovery())
	app.SetTrustedProxies(nil)

	log.Info().Msg("Adding cors, request id and request logging middleware")
	app.Use(middlewares.CORSMiddleware(), middlewares.RequestLogger())

	log.Info().Msg("Setting up routers")
	transport.SetupRouters(app)

	return app
}
