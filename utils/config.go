package utils

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"provider_gateway_mq/consts"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(consts.EnvFile)
	viper.AddConfigPath(consts.EnvFileDirectory)
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Err(err).
			Msg("Error while reading env file, check to OS env config")
	}
	viper.AutomaticEnv()
}

func GetEnvVar(name string) string {
	if !viper.IsSet(name) {
		log.Debug().Msgf("Environment variable %s is not set", name)
		return ""
	}
	value := viper.GetString(name)
	return value
}

func GetCorrelationId() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return "ops" + strconv.FormatInt(t, 10)
}

func GetRmqTlsConf() *tls.Config {
	caCert, err := ioutil.ReadFile(GetEnvVar("_CACERT"))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read CaCert")
	}

	cert, err := tls.LoadX509KeyPair(GetEnvVar("_CERT"), GetEnvVar("_KEY"))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read Certificate, Key pem files")
	}

	rootCA := x509.NewCertPool()
	rootCA.AppendCertsFromPEM(caCert)

	tlsConf := &tls.Config{
		RootCAs:      rootCA,
		Certificates: []tls.Certificate{cert},
	}

	return tlsConf
}
