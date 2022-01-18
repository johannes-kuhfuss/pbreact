package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

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

func setupHttpServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      cfg.RunTime.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func startRouter() {
	var httpsSrv *http.Server
	var m *autocert.Manager

	if cfg.InProduction {
		logger.Info("In production config")
		dataDir := "."
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := cfg.CertDomain
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.TlsPort)
		httpsSrv = setupHttpServer(listenAddr)
		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(dataDir),
			HostPolicy: hostPolicy,
		}

		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				panic(err)
			}
		}()
	}

	listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	httpSrv := setupHttpServer(listenAddr)
	if m != nil {
		httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
	}
	err := httpSrv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
