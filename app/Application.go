package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/handler"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"golang.org/x/crypto/acme/autocert"
)

var (
	cfg config.AppConfig
)

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile, &cfg)
	if err != nil {
		panic(err)
	}
	initRouter()
	mapUrls()
	startRouter()
	logger.Info("Application ended")
}

func initRouter() {
	gin.SetMode(cfg.Gin.Mode)
	gin.DefaultWriter = logger.GetLogger()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)
	cfg.RunTime.Router = router
}

func mapUrls() {
	cfg.RunTime.Router.GET("/ping", handler.Ping)
	cfg.RunTime.Router.GET("/tls", handler.TlsData)
}

func startRouter() {
	err := http.Serve(autocert.NewListener(cfg.CertDomain), cfg.RunTime.Router)
	if err != nil {
		panic(err)
	}
}