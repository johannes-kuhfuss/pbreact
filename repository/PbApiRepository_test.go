package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/stretchr/testify/assert"
)

var (
	cfg  config.AppConfig
	repo PbApiRepository
)

func setupTest(t *testing.T) func() {
	repo = NewPbApiRepository(&cfg)
	return func() {
	}
}

func Test_CreateSubscriptionRequest_Returns_Request(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	jsonTest := dto.PbSubscriptionRequest{}

	req, err := repo.CreateSubscriptionRequest()

	jsonErr := json.Unmarshal(*req, &jsonTest)

	assert.NotNil(t, req)
	assert.Nil(t, err)
	assert.Nil(t, jsonErr)
}

func Test_PrepareHttpRequest_NoRequestType_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiErr := api_error.NewInternalServerError("Request type cannot be empty", nil)
	reqUrl, _ := url.Parse(cfg.PbApi.BaseUrl + "webhooks")

	req, reqErr := repo.PrepareHttpRequest("", reqUrl.String(), nil)

	assert.Nil(t, req)
	assert.NotNil(t, reqErr)
	assert.EqualValues(t, apiErr.StatusCode(), reqErr.StatusCode())
	assert.EqualValues(t, apiErr.Message(), reqErr.Message())
}

func Test_PrepareHttpRequest_WrongRequestType_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiErr := api_error.NewInternalServerError("Could not create http request", nil)
	reqUrl, _ := url.Parse(cfg.PbApi.BaseUrl + "webhooks")

	req, reqErr := repo.PrepareHttpRequest("*?", reqUrl.String(), nil)

	assert.Nil(t, req)
	assert.NotNil(t, reqErr)
	assert.EqualValues(t, apiErr.StatusCode(), reqErr.StatusCode())
	assert.EqualValues(t, apiErr.Message(), reqErr.Message())
}

func Test_PrepareHttpRequest_Returns_Request(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	method := "GET"
	cfg.PbApi.ApiToken = "abcdefg"
	authStr := fmt.Sprintf("Bearer %v", cfg.PbApi.ApiToken)
	reqUrl, _ := url.Parse(cfg.PbApi.BaseUrl + "webhooks")

	req, reqErr := repo.PrepareHttpRequest(method, reqUrl.String(), nil)

	assert.NotNil(t, req)
	assert.Nil(t, reqErr)
	assert.EqualValues(t, method, req.Method)
	assert.EqualValues(t, "1", req.Header.Get("X-Version"))
	assert.EqualValues(t, "application/json", req.Header.Get("Content-Type"))
	assert.EqualValues(t, authStr, req.Header.Get("Authorization"))
}

func Test_ExecHttpRequest_ExecError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	req := &http.Request{}
	apiErr := api_error.NewInternalServerError("Error when executing http request", nil)

	resp, respErr := repo.ExecHttpRequest(req)

	assert.Nil(t, resp)
	assert.NotNil(t, respErr)
	assert.EqualValues(t, apiErr.StatusCode(), respErr.StatusCode())
	assert.EqualValues(t, apiErr.Message(), respErr.Message())
}

func Test_ExecHttpRequest_ErrorStatus_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiErr := api_error.NewInternalServerError("Error when sending request to Productboard API. Status code: 500. Message: Something bad happened", nil)
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something bad happened"))
		}),
	)
	defer srv.Close()
	req, _ := repo.PrepareHttpRequest("GET", srv.URL, nil)

	resp, respErr := repo.ExecHttpRequest(req)

	assert.Nil(t, resp)
	assert.NotNil(t, respErr)
	assert.EqualValues(t, apiErr.StatusCode(), respErr.StatusCode())
	assert.EqualValues(t, apiErr.Message(), respErr.Message())
}

func Test_ExecHttpRequest_Success_Returns_Body(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success"))
		}),
	)
	defer srv.Close()
	req, _ := repo.PrepareHttpRequest("GET", srv.URL, nil)

	resp, respErr := repo.ExecHttpRequest(req)

	assert.NotNil(t, resp)
	assert.Nil(t, respErr)
	assert.EqualValues(t, "Success", string(*resp))
}
