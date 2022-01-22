package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofrs/uuid"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type PbApiServiceInterface interface {
	RegisterForNotifications() api_error.ApiErr
	UnregisterForNotifications() api_error.ApiErr
}

type PbApiService struct {
	Cfg *config.AppConfig
}

func NewPbApiService(cfg *config.AppConfig) PbApiService {
	return PbApiService{
		Cfg: cfg}
}

func (as *PbApiService) RegisterForNotifications() api_error.ApiErr {
	err := as.generateSessionApiToken()
	if err != nil {
		return err
	}
	subReq := dto.PbSubscriptionRequest{
		Data: dto.PbData{
			Name: "Feature Webhooks",
			Events: []dto.Events{
				{EventType: dto.PbEventTypes["featureCreate"]},
				{EventType: dto.PbEventTypes["featureUpdate"]},
				{EventType: dto.PbEventTypes["featureDelete"]},
			},
			Notification: dto.PbNotification{
				URL:     as.Cfg.PbApi.WebHookUrl,
				Version: 1,
				Headers: dto.PbHeaders{
					Authorization: as.Cfg.RunTime.CallbackAuthToken,
				},
			},
		},
	}
	subReqJson, reqErr := json.Marshal(subReq)
	if reqErr != nil {
		msg := "Could not generate subscription request"
		logger.Error(msg, err)
		return api_error.NewInternalServerError(msg, err)
	}
	subscriptionUrl, _ := url.Parse(as.Cfg.PbApi.BaseUrl + "webhooks")
	req, reqErr := http.NewRequest("POST", subscriptionUrl.String(), bytes.NewBuffer(subReqJson))
	if reqErr != nil {
		msg := "Could not create subscription http request"
		logger.Error(msg, reqErr)
		return api_error.NewInternalServerError(msg, reqErr)
	}
	authStr := fmt.Sprintf("Bearer %v", as.Cfg.PbApi.ApiToken)
	req.Header = http.Header{
		"X-Version":     []string{"1"},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{authStr},
	}
	client := http.Client{}
	resp, resErr := client.Do(req)
	if resErr != nil {
		msg := "Error when trying to subscribe for notifications"
		logger.Error(msg, reqErr)
		return api_error.NewInternalServerError(msg, reqErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		msg := fmt.Sprintf("Error when sending registration request to API. Status code: %v. Message: %v", resp.StatusCode, string(body))
		logger.Error(msg, nil)
		return api_error.NewInternalServerError(msg, nil)
	} else {
		logger.Info(fmt.Sprintf("Successfully registered for notifications. Status code: %v", resp.StatusCode))
	}
	return nil
}

func (as *PbApiService) generateSessionApiToken() api_error.ApiErr {
	logger.Info("Generating API callback token")
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
	notifs, err := as.GetNotifications()
	if err != nil {
		return err
	}
	for _, val := range notifs.Data {
		id := val.ID
		subscriptionUrl, _ := url.Parse(as.Cfg.PbApi.BaseUrl + "webhooks" + "/" + id)
		req, reqErr := http.NewRequest("DELETE", subscriptionUrl.String(), nil)
		if reqErr != nil {
			msg := "Could not create delete subscription http request"
			logger.Error(msg, reqErr)
			return api_error.NewInternalServerError(msg, reqErr)
		}
		authStr := fmt.Sprintf("Bearer %v", as.Cfg.PbApi.ApiToken)
		req.Header = http.Header{
			"X-Version":     []string{"1"},
			"Authorization": []string{authStr},
		}
		client := http.Client{}
		resp, resErr := client.Do(req)
		if resErr != nil {
			msg := "Error when trying to delete subscriptions for notifications"
			logger.Error(msg, reqErr)
			return api_error.NewInternalServerError(msg, reqErr)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode > 299 {
			msg := fmt.Sprintf("Error when sending delete subscriptions request to API. Status code: %v. Message: %v", resp.StatusCode, string(body))
			logger.Error(msg, nil)
			return api_error.NewInternalServerError(msg, nil)
		} else {
			logger.Info(fmt.Sprintf("Successfully deleted subscriptions for notifications. Status code: %v", resp.StatusCode))
		}
	}
	return nil
}

func (as *PbApiService) GetNotifications() (*dto.PbSubscriptionResponse, api_error.ApiErr) {
	var pbResp dto.PbSubscriptionResponse
	subscriptionUrl, _ := url.Parse(as.Cfg.PbApi.BaseUrl + "webhooks")
	req, reqErr := http.NewRequest("GET", subscriptionUrl.String(), nil)
	if reqErr != nil {
		msg := "Could not create list subscription http request"
		logger.Error(msg, reqErr)
		return nil, api_error.NewInternalServerError(msg, reqErr)
	}
	authStr := fmt.Sprintf("Bearer %v", as.Cfg.PbApi.ApiToken)
	req.Header = http.Header{
		"X-Version":     []string{"1"},
		"Authorization": []string{authStr},
	}
	client := http.Client{}
	resp, resErr := client.Do(req)
	if resErr != nil {
		msg := "Error when trying to list subscriptions for notifications"
		logger.Error(msg, reqErr)
		return nil, api_error.NewInternalServerError(msg, reqErr)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		msg := fmt.Sprintf("Error when sending list subscriptions request to API. Status code: %v. Message: %v", resp.StatusCode, string(body))
		logger.Error(msg, nil)
		return nil, api_error.NewInternalServerError(msg, nil)
	} else {
		logger.Info(fmt.Sprintf("Successfully listed subscriptions for notifications. Status code: %v", resp.StatusCode))
	}

	err := json.Unmarshal(body, &pbResp)
	if err != nil {
		msg := "Error parsing subscription list"
		logger.Error(msg, err)
		return nil, api_error.NewInternalServerError(msg, err)
	}
	if len(pbResp.Data) == 0 {
		msg := "No subscriptions found"
		logger.Error(msg, nil)
		return nil, api_error.NewNotFoundError(msg)
	}
	return &pbResp, nil
}
