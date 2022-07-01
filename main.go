package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"provider_gateway_mq/app"
	"provider_gateway_mq/controllers"
	"provider_gateway_mq/utils"
)

func init() {
	// Set gin mode
	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	// Setup the app gateway

	host := utils.GetEnvVar("GIN_ADDR")
	port := utils.GetEnvVar("GIN_PORT")
	https := utils.GetEnvVar("GIN_HTTPS")
	connectionString := utils.GetEnvVar("RMQ_URL")

	configs := utils.GetQueueConf()

	rmqProducer := controllers.RMQSpec{
		ConnectionString: connectionString,
		Err:              make(chan error),
	}

	rmqConsumer := controllers.RMQSpec{
		ConnectionString: connectionString,
		Err:              make(chan error),
	}

	err := rmqProducer.Connect()
	if err != nil {
		rmqProducer.OnError(err, "Failed to connect to RabbitMQ while publishing")
	}

	err = rmqConsumer.Connect()
	if err != nil {
		rmqConsumer.OnError(err, "Failed to declare a queue while consuming")
	}

	for _, conf := range configs {

		rmqProducer.Exchange = conf.Topic
		rmqProducer.Queue = conf.QueueName
		rmqProducer.BindingKey = conf.BindingKey
		rmqProducer.RoutingKey = conf.RoutingKey
		rmqProducer.ReplyTo = conf.ReplyTo

		rmqConsumer.Queue = conf.ReplyTo

		rmqProducer.PublishDeclare()
		rmqConsumer.ConsumeDeclare()
	}

	go rmqProducer.PublishMessages()
	go rmqConsumer.ConsumeMessages()

	appApiGateway := app.SetupApp()

	if https == "true" {
		certFile := utils.GetEnvVar("GIN_CERT")
		certKey := utils.GetEnvVar("GIN_CERT_KEY")
		log.Info().Msgf("Starting service on https://%s:%s", host, port)

		if err := appApiGateway.RunTLS(fmt.Sprintf("%s:%s", host, port), certFile, certKey); err != nil {
			log.Fatal().Err(err).Msg("Error on setting up the server in HTTPS mode")
		}
	}

	log.Info().Msgf("Starting service on http://%s:%s", host, port)
	if err := appApiGateway.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Fatal().Err(err).Msg("Error on setting up the server")
	}
}
