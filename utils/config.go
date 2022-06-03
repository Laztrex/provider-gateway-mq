package utils

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"provider_gateway_mq/consts"
)

func init() {
	viper.SetConfigFile(consts.EnvFile)
	viper.AddConfigPath(consts.EnvFileDirectory)
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Err(err).
			Msg("Error occurred while reading env file, might fallback to OS env config")
	}
	viper.AutomaticEnv()
}

// GetEnvVar This function can be used to get ENV Var in our App
// Modify this if you change the library to read ENV
func GetEnvVar(name string) string {
	if !viper.IsSet(name) {
		log.Debug().Msgf("Environment variable %s is not set", name)
		return ""
	}
	value := viper.GetString(name)
	return value
}

func GetTlsConf() *tls.Config {
	caCert, err := ioutil.ReadFile(GetEnvVar("_CACERT"))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read CACert")
	}

	cert, err := tls.LoadX509KeyPair(GetEnvVar("_CERT"), GetEnvVar("_KEY"))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read Certificate, Key")
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)

	tlsConf := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
	}

	return tlsConf
}

func GetCorrelationId() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return "msg" + strconv.FormatInt(t, 10)
}

func Timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
