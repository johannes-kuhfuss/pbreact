package app

import (
	"fmt"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/handler"
	"github.com/johannes-kuhfuss/services_utils/logger"
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
	getCert()
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
}

func getCert() {
	err := autotls.Run(cfg.RunTime.Router, cfg.CertDomain)
	if err != nil {
		logger.Error("Error while getting certificate", err)
		panic(err)
	}
}

func startRouter() {
	listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))

	if err := cfg.RunTime.Router.Run(listenAddr); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}
