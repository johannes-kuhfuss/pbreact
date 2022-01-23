package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/mocks/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/stretchr/testify/assert"
)

var (
	cfg    config.AppConfig
	whh    WebHookHandler
	router *gin.Engine
	//mockService *service.MockPbApiService
	recorder *httptest.ResponseRecorder
	ctx      *gin.Context
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
	apiError := api_error.NewUnauthenticatedError("Wrong or missing auth key")
	errorJson, _ := json.Marshal(apiError)
	router.GET("/pbwebhook", whh.PbWhSubscription)
	req, _ := http.NewRequest(http.MethodGet, "/pbwebhook", nil)

	router.ServeHTTP(recorder, req)

	assert.EqualValues(t, apiError.StatusCode(), recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_PbWhSubscription_NoId_Returns_BadRequestError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("Could not find validation token")
	errorJson, _ := json.Marshal(apiError)
	authKey, _ := uuid.NewV4()
	cfg.RunTime.CallbackAuthToken = authKey.String()
	router.GET("/pbwebhook", whh.PbWhSubscription)
	req, _ := http.NewRequest(http.MethodGet, "/pbwebhook", nil)
	req.Header.Set("Authorization", authKey.String())

	router.ServeHTTP(recorder, req)

	assert.EqualValues(t, apiError.StatusCode(), recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_PbWhSubscription_WithId_Returns_Id(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	authKey, _ := uuid.NewV4()
	cfg.RunTime.CallbackAuthToken = authKey.String()
	router.GET("/pbwebhook", whh.PbWhSubscription)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/pbwebhook?validationToken=%v", authKey.String()), nil)
	req.Header.Set("Authorization", authKey.String())

	router.ServeHTTP(recorder, req)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, authKey.String(), recorder.Body.String())
}

func Test_PbWhEvents_NoAuthKey_Returns_UnauthenticatedError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewUnauthenticatedError("Wrong or missing auth key")
	errorJson, _ := json.Marshal(apiError)
	router.POST("/pbwebhook", whh.PbWhEvents)
	req, _ := http.NewRequest(http.MethodPost, "/pbwebhook", nil)

	router.ServeHTTP(recorder, req)

	assert.EqualValues(t, apiError.StatusCode(), recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_PbWhEvents_InvalidBody_Returns_BadRequestError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("Invalid json body")
	errorJson, _ := json.Marshal(apiError)
	authKey, _ := uuid.NewV4()
	cfg.RunTime.CallbackAuthToken = authKey.String()
	router.POST("/pbwebhook", whh.PbWhEvents)
	req, _ := http.NewRequest(http.MethodPost, "/pbwebhook", nil)
	req.Header.Set("Authorization", authKey.String())

	router.ServeHTTP(recorder, req)

	assert.EqualValues(t, apiError.StatusCode(), recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}
