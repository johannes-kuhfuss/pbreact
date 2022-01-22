package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/domain"
	"github.com/johannes-kuhfuss/pbreact/handler"
	"github.com/johannes-kuhfuss/pbreact/repository"
	"github.com/johannes-kuhfuss/pbreact/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	cfg          config.AppConfig
	pbApiRepo    domain.PbApiRepository
	pbApiService service.DefaultPbApiService
	pbApiHandler handler.WebHookHandler
	server       http.Server
	appEnd       chan os.Signal
	ctx          context.Context
	cancel       context.CancelFunc
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
	RegisterForOsSignals()
	go RegisterForNotifications()
	go startServer()

	<-appEnd
	cleanUp()

	if srvErr := server.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", srvErr)
	} else {
		logger.Info("Graceful shutdown finished")
	}
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
	}
}

func wireApp() {
	pbApiRepo = repository.NewPbApiRepository(&cfg)
	pbApiService = service.NewPbApiService(&cfg, pbApiRepo)
	pbApiHandler = handler.NewWebHookHandler(&cfg, pbApiService)
}

func mapUrls() {
	cfg.RunTime.Router.GET("/ping", handler.Ping)
	cfg.RunTime.Router.GET("/pbwebhook", pbApiHandler.PbWhSubscription)
	cfg.RunTime.Router.POST("/pbwebhook", pbApiHandler.PbWhEvents)
}

func RegisterForOsSignals() {
	appEnd = make(chan os.Signal, 1)
	signal.Notify(appEnd, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func RegisterForNotifications() {
	time.Sleep(10 * time.Second)
	logger.Info("Registering for notifications")
	err := pbApiService.RegisterForNotifications()
	if err != nil {
		panic(err)
	}
}

func startServer() {
	listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.TlsPort)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	if err := server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil && err != http.ErrServerClosed {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}

func cleanUp() {
	shutdownTime := time.Duration(cfg.GracefulShutdownTime) * time.Second
	ctx, cancel = context.WithTimeout(context.Background(), shutdownTime)
	defer func() {
		logger.Info("Cleaning up")
		pbApiService.UnregisterForNotifications()
		logger.Info("Done cleaning up")
		cancel()
	}()
}
