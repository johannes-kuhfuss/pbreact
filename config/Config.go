package config

import (
	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Server struct {
		Host     string `envconfig:"SERVER_HOST"`
		Port     string `envconfig:"SERVER_PORT" default:"8080"`
		TlsPort  string `envconfig:"SERVER_TLSPORT" default:"8443"`
		Shutdown bool   `ignored:"true" default:"false"`
		CertFile string `envconfig:"CERT_FILE" default:"./cert/cert.pem"`
		KeyFile  string `envconfig:"KEY_FILE" default:"./cert/cert.key"`
	}
	Gin struct {
		Mode string `envconfig:"GIN_MODE" default:"release"`
	}
	RunTime struct {
		Router *gin.Engine
	}
	PbAuthHeader string `envconfig:"PB_AUTH_HEADER"`
}

const (
	EnvFile = ".env"
)

func InitConfig(file string, config *AppConfig) api_error.ApiErr {
	logger.Info("Initalizing configuration")
	loadConfig(file)
	err := envconfig.Process("", config)
	if err != nil {
		return api_error.NewInternalServerError("Could not initalize configuration. Check your environment variables", err)
	}
	logger.Info("Done initalizing configuration")
	return nil
}

func loadConfig(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		logger.Error("Could not open env file", err)
		return err
	}
	return nil
}
