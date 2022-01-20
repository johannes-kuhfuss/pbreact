package service

import (
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type PbApiServiceInterface interface {
	RegisterForNotifications() api_error.ApiErr
}

type PbApiService struct {
	Cfg *config.AppConfig
}

func NewPbApiService(cfg *config.AppConfig) PbApiService {
	return PbApiService{cfg}
}

func (as *PbApiService) RegisterForNotifications() api_error.ApiErr {
	return nil
}
