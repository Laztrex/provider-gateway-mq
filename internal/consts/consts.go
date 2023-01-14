package consts

const EnvFile = ".env"
const EnvFileDirectory = "."

const QueuesConf = "configs/queue_config.yaml"

const (
	KeyRequestId     = "RqUID"
	KeyCorrelationId = "CorrID"
	KeyRoutingKey    = "routing-key"
)

const LogPath = "/var/log/metrics.log"

const (
	MqCaCertDefault = "/certs/gateway_mq/cacert.pem"
	MqCertDefault   = "/certs/gateway_mq/client_cert.pem"
	MqKeyDefault    = "/certs/gateway_mq/client_key.pem"
)

const (
	GinCertDefault    = "/certs/gateway_mq/server_cert.pem"
	GinCertKeyDefault = "/certs/gateway_mq/server_key.pem"
)
