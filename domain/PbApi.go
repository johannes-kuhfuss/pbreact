package domain

import (
	"github.com/johannes-kuhfuss/pbreact/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type PbApiRepository interface {
	RegisterForNotifications() api_error.ApiErr
	GetNotifications() (*dto.PbSubscriptionResponse, api_error.ApiErr)
	UnregisterForNotifications(dto.PbSubscriptionResponse) api_error.ApiErr
}
