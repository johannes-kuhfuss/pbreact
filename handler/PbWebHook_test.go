package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/mocks/service"
)

var (
	cfg         config.AppConfig
	whh         WebHookHandler
	router      *gin.Engine
	mockService *service.MockPbApiService
	recorder    *httptest.ResponseRecorder
)

func setupTest(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockPbApiService(ctrl)
	whh = NewWebHookHandler(&cfg, mockService)
	router = gin.Default()
	recorder = httptest.NewRecorder()
	return func() {
		router = nil
		ctrl.Finish()
	}
}
