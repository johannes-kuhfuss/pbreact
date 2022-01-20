package app

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/handler"
	"github.com/johannes-kuhfuss/pbreact/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	cfg          config.AppConfig
	whh          handler.WebHookHandler
	pbApiService service.PbApiService
	server       http.Server
)

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile, &cfg)
	if err != nil {
		panic(err)
	}
	initRouter()
	initServer()
	wireApp()
	mapUrls()
	go startServer()
	registerForNotif()
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

func initServer() {
	tlsConfig := tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	}
	listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.TlsPort)
	server = http.Server{
		Addr:              listenAddr,
		Handler:           cfg.RunTime.Router,
		TLSConfig:         &tlsConfig,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    0,
		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		//ErrorLog:          &log.Logger{},
	}
}

func wireApp() {
	whh = handler.WebHookHandler{
		Cfg: &cfg,
	}
	pbApiService = service.NewPbApiService(&cfg)
}

func mapUrls() {
	cfg.RunTime.Router.GET("/ping", handler.Ping)
	cfg.RunTime.Router.GET("/pbwebhook", whh.PbWhSubscription)
	cfg.RunTime.Router.POST("/pbwebhook", whh.PbWhEvents)
}

func startServer() {
	listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.TlsPort)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	if err := server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}

func registerForNotif() {
	logger.Info("Registering for notifications")
	err := pbApiService.RegisterForNotifications()
	if err != nil {
		logger.Error("Error while registering for notifications", err)
	}
}
