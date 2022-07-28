package consts

const EnvFile = ".env"
const EnvFileDirectory = "."

const QueuesConf = "configs/queue_config.yaml"

const RequestIdHttpHeaderName = "request-id"

const LogPath = "/var/log/metrics.log"

const (
	MqCACERT = "/certs/gateway_mq/cacert_mq.pem"
	MqCERT   = "/certs/gateway_mq/client_cert_mq.pem"
	MqKEY    = "/certs/gateway_mq/client_key_mq.pem"
)

const (
	GinCert    = "/certs/gateway_mq/server_cert.pem"
	GinCertKey = "/certs/gateway_mq/server_key.pem"
)
