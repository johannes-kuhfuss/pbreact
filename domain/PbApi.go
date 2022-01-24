package domain

import (
	"github.com/johannes-kuhfuss/pbreact/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

//go:generate mockgen -destination=../mocks/domain/mockPbApiRepository.go -package=domain github.com/johannes-kuhfuss/pbreact/domain PbApiRepository
type PbApiRepository interface {
	RegisterForNotifications() api_error.ApiErr
	GetNotifications() (*dto.PbSubscriptionResponse, api_error.ApiErr)
	UnregisterForNotifications(dto.PbSubscriptionResponse) api_error.ApiErr
}
