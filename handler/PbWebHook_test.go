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
	ctx         *gin.Context
)

func setupTest(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockPbApiService(ctrl)
	whh = NewWebHookHandler(&cfg, mockService)
	router = gin.Default()
	gin.SetMode(gin.TestMode)
	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	return func() {
		router = nil
		ctrl.Finish()
	}
}

func Test_PbWhSubscription_NoAuthKey_Returns_UnauthenticatedError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
}
