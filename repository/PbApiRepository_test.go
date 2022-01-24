package repository

import (
	"encoding/json"
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

	req, reqErr := repo.PrepareHttpRequest("", "", nil)

	assert.Nil(t, req)
	assert.NotNil(t, reqErr)
	assert.EqualValues(t, apiErr.StatusCode(), reqErr.StatusCode())
	assert.EqualValues(t, apiErr.Message(), reqErr.Message())
}
