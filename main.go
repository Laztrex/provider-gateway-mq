package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"provider_gateway_mq/app"
	"provider_gateway_mq/consts"
	"provider_gateway_mq/controllers"
	"provider_gateway_mq/utils"
)

func init() {
	// Set gin mode

	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	connectionString := utils.GetEnvVar("RMQ_URL")

	rmqProducer := controllers.RMQProducer{
		Queue:            consts.QueueName,
		ConnectionString: connectionString,
	}

	rmqConsumer := controllers.RMQConsumer{
		Queue:            consts.AnswerQueueName,
		ConnectionString: connectionString,
	}

	go rmqProducer.PublishMessages()
	go rmqConsumer.ConsumeMessages()

	appApiGateway := app.SetupApp()

	host := utils.GetEnvVar("GIN_ADDR")
	port := utils.GetEnvVar("GIN_PORT")

	https := utils.GetEnvVar("GIN_HTTPS")
	if https == "true" {
		certFile := utils.GetEnvVar("GIN_CERT")
		certKey := utils.GetEnvVar("GIN_CERT_KEY")
		log.Info().Msgf("Starting service on https://%s:%s", host, port)

		if err := appApiGateway.RunTLS(fmt.Sprintf("%s:%s", host, port), certFile, certKey); err != nil {
			log.Fatal().Err(err).Msg("Error while setting up the server in HTTPS mode")
		}
	}
	log.Info().Msgf("Starting service on http://%s:%s", host, port)
	if err := appApiGateway.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Fatal().Err(err).Msg("Error while setting up the server")
	}

}
