package service

import (
	"github.com/gofrs/uuid"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type PbApiServiceInterface interface {
	RegisterForNotifications() api_error.ApiErr
	UnregisterForNotifications() api_error.ApiErr
}

type PbApiService struct {
	repo domain.PbApiRepository
	Cfg  *config.AppConfig
}

func NewPbApiService(cfg *config.AppConfig, repo domain.PbApiRepository) PbApiService {
	return PbApiService{
		repo: repo,
		Cfg:  cfg,
	}
}

func (as *PbApiService) RegisterForNotifications() api_error.ApiErr {
	var err api_error.ApiErr
	err = as.generateSessionApiToken()
	if err != nil {
		return err
	}
	err = as.repo.RegisterForNotifications()
	if err != nil {
		return err
	}
	return nil
}

func (as *PbApiService) generateSessionApiToken() api_error.ApiErr {
	id, err := uuid.NewV4()
	if err != nil {
		msg := "Could not generate callback auth token"
		logger.Error(msg, err)
		return api_error.NewInternalServerError(msg, err)
	}
	as.Cfg.RunTime.CallbackAuthToken = id.String()
	return nil
}

func (as *PbApiService) UnregisterForNotifications() api_error.ApiErr {
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
