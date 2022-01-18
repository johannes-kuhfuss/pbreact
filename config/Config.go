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
		Port     string `envconfig:"SERVER_PORT" default:"80"`
		TlsPort  string `envconfig:"SERVER_TLSPORT" default:"443"`
		Shutdown bool   `ignored:"true" default:"false"`
	}
	CertDomain   string `envconfig:"CERT_DOMAIN"`
	InProduction bool   `envconfig:"IN_PRODUCTION" default:"true"`
	Gin          struct {
		Mode string `envconfig:"GIN_MODE" default:"release"`
	}
	RunTime struct {
		Router *gin.Engine
	}
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
