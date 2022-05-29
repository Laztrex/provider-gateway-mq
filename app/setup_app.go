package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"provider_gateway_mq/middlewares"
	"provider_gateway_mq/routers"
)

func SetupApp() *gin.Engine {
	log.Info().Msg("Initializing service")

	// Create engine
	app := gin.New()

	app.Use(gin.Recovery())

	app.SetTrustedProxies(nil)

	log.Info().Msg("Adding cors, request_id and request logging middlewares")
	app.Use(middlewares.CORSMiddleware(), middlewares.RequestID(), middlewares.RequestLogger())

	log.Info().Msg("Setting up routers")
	routers.SetupRouters(app)

	return app

}
