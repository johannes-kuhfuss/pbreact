package service

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/dto"
	"github.com/johannes-kuhfuss/pbreact/mocks/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/stretchr/testify/assert"
)

var (
	pbApiCtrl     *gomock.Controller
	mockPbApiRepo *domain.MockPbApiRepository
	as            PbApiService
	cfg           config.AppConfig
)

func setupApi(t *testing.T) func() {
	pbApiCtrl = gomock.NewController(t)
	mockPbApiRepo = domain.NewMockPbApiRepository(pbApiCtrl)
	as = NewPbApiService(&cfg, mockPbApiRepo)
	return func() {
		as = nil
		pbApiCtrl.Finish()
	}
}

func isValidUuid(id string) bool {
	_, err := uuid.FromString(id)
	return err == nil
}

func Test_GenerateSessionApiToken_Returns_ApiToken(t *testing.T) {
	teardown := setupApi(t)
	defer teardown()

	err := as.GenerateSessionApiToken()

	assert.Nil(t, err)
	id := cfg.RunTime.CallbackAuthToken
	assert.EqualValues(t, true, isValidUuid(id))
}

func Test_RegisterForNotifications_Returns_Error(t *testing.T) {
	teardown := setupApi(t)
	defer teardown()
	apiError := api_error.NewInternalServerError("something went wrong", nil)

	mockPbApiRepo.EXPECT().RegisterForNotifications().Return(apiError)

	err := as.RegisterForNotifications()

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_RegisterForNotifications_Returns_NoError(t *testing.T) {
	teardown := setupApi(t)
	defer teardown()

	mockPbApiRepo.EXPECT().RegisterForNotifications().Return(nil)

	err := as.RegisterForNotifications()

	assert.Nil(t, err)
}

func Test_UnregisterForNotifications_GetFails_Returns_Error(t *testing.T) {
	teardown := setupApi(t)
	defer teardown()
	apiError := api_error.NewInternalServerError("something went wrong", nil)

	mockPbApiRepo.EXPECT().GetNotifications().Return(nil, apiError)

	err := as.UnregisterForNotifications()

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_UnregisterForNotifications_UnregisterFails_Returns_Error(t *testing.T) {
	teardown := setupApi(t)
	defer teardown()
	apiError := api_error.NewInternalServerError("something went wrong", nil)
	notifs := dto.PbSubscriptionResponse{
		Data:  []dto.SubRespData{},
		Links: dto.Links{},
	}

	mockPbApiRepo.EXPECT().GetNotifications().Return(&notifs, nil)
	mockPbApiRepo.EXPECT().UnregisterForNotifications(notifs).Return(apiError)

	err := as.UnregisterForNotifications()

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_UnregisterForNotifications_Returns_NoError(t *testing.T) {
	teardown := setupApi(t)
	defer teardown()
	notifs := dto.PbSubscriptionResponse{
		Data:  []dto.SubRespData{},
		Links: dto.Links{},
	}

	mockPbApiRepo.EXPECT().GetNotifications().Return(&notifs, nil)
	mockPbApiRepo.EXPECT().UnregisterForNotifications(notifs).Return(nil)

	err := as.UnregisterForNotifications()

	assert.Nil(t, err)
}
