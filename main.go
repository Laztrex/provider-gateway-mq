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
	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {

	host := utils.GetEnvVar("GIN_ADDR")
	port := utils.GetEnvVar("GIN_PORT")
	https := utils.GetEnvVar("GIN_HTTPS")

	connectionString := utils.GetEnvVar("RMQ_URL")
	exchange := utils.GetEnvVar("TOPIC")
	routingKey := utils.GetEnvVar("ROUTING_KEY")
	bindingKey := utils.GetEnvVar("BINDING_KEY")
	queue := utils.GetEnvVar("QUEUE_NAME")

	rmqProducer := controllers.RMQSpec{
		Queue:            queue,
		ConnectionString: connectionString,
		Exchange:         exchange,
		RoutingKey:       routingKey,
		BindingKey:       bindingKey,
		Err:              make(chan error),
	}

	//rmqProducer2 := controllers.RMQSpec{
	//	Queue:            "ml361",
	//	ConnectionString: connectionString,
	//	Exchange:         exchange,
	//	RoutingKey:       routingKey,
	//	BindingKey:       "*.metrics",
	//	Err:              make(chan error),
	//}
	//rmqProducer2.PublishConnecting()

	rmqConsumer := controllers.RMQSpec{
		Queue:            consts.AnswerQueueName,
		ConnectionString: connectionString,
		Exchange:         exchange,
		RoutingKey:       routingKey,
		Err:              make(chan error),
	}

	rmqProducer.PublishConnecting()

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
