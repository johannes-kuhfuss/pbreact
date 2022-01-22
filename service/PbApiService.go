package service

import (
	"github.com/gofrs/uuid"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

//go:generate mockgen -destination=../mocks/service/mockPbApiService.go -package=service github.com/johannes-kuhfuss/pbreact/service PbApiService
type PbApiService interface {
	RegisterForNotifications() api_error.ApiErr
	UnregisterForNotifications() api_error.ApiErr
}

type DefaultPbApiService struct {
	repo domain.PbApiRepository
	cfg  *config.AppConfig
}

func NewPbApiService(c *config.AppConfig, r domain.PbApiRepository) DefaultPbApiService {
	return DefaultPbApiService{
		repo: r,
		cfg:  c,
	}
}

func (as DefaultPbApiService) RegisterForNotifications() api_error.ApiErr {
	err := as.generateSessionApiToken()
	if err != nil {
		return err
	}
	err = as.repo.RegisterForNotifications()
	if err != nil {
		return err
	}
	return nil
}

func (as DefaultPbApiService) generateSessionApiToken() api_error.ApiErr {
	id, err := uuid.NewV4()
	if err != nil {
		msg := "Could not generate callback auth token"
		logger.Error(msg, err)
		return api_error.NewInternalServerError(msg, err)
	}
	as.cfg.RunTime.CallbackAuthToken = id.String()
	return nil
}

func (as DefaultPbApiService) UnregisterForNotifications() api_error.ApiErr {
	notifs, err := as.repo.GetNotifications()
	if err != nil {
		return err
	}
	err = as.repo.UnregisterForNotifications(*notifs)
	if err != nil {
		return err
	}
	return nil
}
