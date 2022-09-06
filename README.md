# Provider-gateway-mq

The project is part of a cloud-based microservice solution for interacting with an ML-application and related
infrastructure components.

This application allows the Initiator to communicate with the model through the RabbitMQ client in the cluster
project.    
*Consumer* can be located both inside the interface of the model, and in a dedicated application that interacts with the
model using the RPC-pattern.

## Description

*Provider-gateway-mq* currently implements the following:

* REST request
* Application arranges sync/async RPC calls with the model
* The routing path to publish to the queue must be contained in the "routing_key" header. Also, the gateway-mq configuration settings provide for setting a fixed value of "routing_key"
* When the model needs to transmit the result of its work, it turns to another (Reply) queue to transmit the gateway-mq response and sends a message with the result of the work there
* Response from the model (sync)
    * If the pod with the model crashes before being called or during the calculation - the message will not be lost
      under certain conditions
    * In the event of a timeout from the side of the Initiator, provide for the possibility of saving the response and
      issuing "uid" if the response did come (*in progress*)
* Declare of exchange, queues on the *provider-gateway-mq* side - we do not trust the model

The project will be updated.

Tested on
~~~
>> go version
go version go1.17.2 darwin/amd64
~~~
~~~
>> rabbitmqctl version
3.10.6
~~~

## Usage

The project uses global *environment variables*

| Env Name     |           Goal           |         Expected value example         |
|--------------|:------------------------:|:--------------------------------------:|
| RMQ_URL      | Host to connect RabbitMQ | "amqps://user:password@rabbitmq:5671/" |
| GIN_MODE     |          label           |               "release"                |
| GIN_HTTPS    |         TLS mode         |                 "true"                 |
| GIN_ADDR     |           Host           |               "0.0.0.0"                |
| GIN_PORT     |           Port           |                 "5050"                 |
| LOG_LEVEL    |      Logging level       |                   "debug"                   |

and *constants*.

| Const Name              |                Goal                |         Expected value example         |
|-------------------------|:----------------------------------:|:--------------------------------------:|
| MqCACERT                |                 CA                 |   "/certs/gateway_mq/cacert_mq.pem"    |
| MqCERT                  |            Certificate             | "/certs/gateway_mq/client_cert_mq.pem" |
| MqKEY                   |          Key certificate           | "/certs/gateway_mq/client_key_mq.pem"  |
| GinCert                 |         Server certificate         |  "/certs/gateway_mq/server_cert.pem"   |
| GinCertKey              |             Server key             |   "/certs/gateway_mq/server_key.pem"   |
| EnvFile                 |     Local env file in project      | ".env"  |
| EnvFileDirectory        |            Dir env file            |  "."   |
| QueuesConf              | Path to configration Queue declare |   "configs/queue_config.yaml"   |
| RequestIdHttpHeaderName |    Name Header for *request-id*    | ".env"  |
| LogPath                 |         Path to dump logs          |   "/var/log/metrics.log"   |

You can define the default values in the **.env** file of the project root.

Configuration file is also provided to define protocol settings
- [queue_config.yml](https://github.com/Laztrex/provider-gateway-mq/blob/main/configs/)

~~~yaml
- topic: ML.MQ
  queueName: fib
  bindingKey: "predict.*"
  replyTo: "response"
~~~

Currently, only Topic can be defined from the configuration file for Exchange (but Direct can also be defined). Other
exchanges - will be supplemented. Flexible settings for the queue - lifetime, autodelete, types, arguments and etc.
currently not included in the configuration file, the parameters can be configured inside the code optionally.

The directory [examples/webapp](https://github.com/Laztrex/provider-gateway-mq/blob/main/examples/webapp/) contains a
simple web application for testing the project.

~~~
>> go version
go version go1.18.2 darwin/amd64
~~~

### Compose it

The project provides examples of services for testing provider-gateway-mq.

~~~
docker-compose build --no-cache
~~~

~~~
docker-compose up -d
~~~

~~~
docker-compose logs -f -t
~~~

Check:

~~~
>> examples % curl -k --key certs/client/client_key.pem --cert certs/client/client_cert.pem https://127.0.0.1:5050/v1/predict -d '{"data": "[15, 29]"}' -H "RqUID: 52-42" -H "Content-Type: application/json" -H "routing-key: predict.online"
~~~

## Objective

The goal of the project is to create binding services for the ML model that provide ease of interaction, flexible logic,
scalability of connecting various interfaces of the ML Engine architecture - *Cloud-based Sandbox for ML Serving*

Simplified sketch  
![Image alt](https://github.com/Laztrex/provider-gateway-mq/blob/main/docs/pics/first_sketch.png)

## Addition

It is implied in the microservice architecture of a cloud solution for calculating ML models:

* Prefix "**provider**-x-x" - auxiliary controllers to implement the interface with the infrastructure technical stack
* Prefix "x-**gateway**-x" - optional label to identify the converter, here indicates the presence of a *REST-MQ*
  adapter
* Predix "x-x-**mq**" - integrated interface

### Certificates

In this project, we are going to create a Golang web client to connect to RabbitMQ server with TLS. For this we you will
need to create self-signed SSL certificates and share them between the Golang application and the RabbitMQ server.

Directory [examples/certs](https://github.com/Laztrex/provider-gateway-mq/blob/main/examples/certs/) contains a
Dockerfile with an example of generating a self-signed certificate.

~~~
cd examples/certs/
~~~

~~~
docker build -t certs .
~~~

~~~
docker run -i -t certs bash
>> cd /tls-gen/basic/result
~~~

Add certificate files to
directory [examples/certs](https://github.com/Laztrex/provider-gateway-mq/blob/main/examples/certs/). Initially set
certificate structure:

~~~
├── examples  
│   └── certs 
│       ├── rabbitmq
│       │    ├── cacert.pem  
│       │    ├── server_cert.pem  
│       │    └── server_key.pem  
│       ├── gateway_mq
│       │    ├── cacert.pem  
│       │    ├── client_cert.pem  
│       │    └── client_key.pem 
│       └── webapp
│            ├── cacert.pem  
│            ├── client_cert.pem  
│            └── client_key.pem   
└── docker-compose.yaml  
~~~

It is also necessary to provide for the creation of certificates for the *gin*-server.

### Initial setup RabbitMQ

When initializing the RabbitMQ client, you can set initial settings.  
In particular, you can set available accounts and define a configuration file. See example
in [examples/rabbit-mq](https://github.com/Laztrex/provider-gateway-mq/blob/main/examples/rabbit-mq/).

Check the list of available users:

~~~
>> rabbitmqctl list_users
superuser	[administrator]
MLUser	[consumer]
~~~

Check declared queues:

~~~
>> rabbitmqctl list_queues
name	messages
response	1
manage	0
fib	1
~~~