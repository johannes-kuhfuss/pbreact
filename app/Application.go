package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/handler"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	cfg config.AppConfig
	whh handler.WebHookHandler
)

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile, &cfg)
	if err != nil {
		panic(err)
	}
	initRouter()
	wireApp()
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

func wireApp() {
	whh = handler.WebHookHandler{
		Cfg: &cfg,
	}
}

func mapUrls() {
	cfg.RunTime.Router.GET("/ping", handler.Ping)
	cfg.RunTime.Router.GET("/pbwebhook", whh.PbWebHook)
}

func startRouter() {
	listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.TlsPort)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	if err := cfg.RunTime.Router.RunTLS(listenAddr, cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}
