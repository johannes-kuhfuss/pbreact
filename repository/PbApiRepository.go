package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type PbApiRepository struct {
	cfg *config.AppConfig
}

func NewPbApiRepository(c *config.AppConfig) PbApiRepository {
	return PbApiRepository{
		cfg: c,
	}
}

func (r PbApiRepository) RegisterForNotifications() api_error.ApiErr {
	subReq, err := r.CreateSubscriptionRequest()
	if err != nil {
		return err
	}
	reqUrl, _ := url.Parse(r.cfg.PbApi.BaseUrl)
	reqUrl.Path = "/webhook"
	req, err := r.PrepareHttpRequest("POST", reqUrl.String(), bytes.NewBuffer(*subReq))
	if err != nil {
		return err
	}
	_, err = r.ExecHttpRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (r PbApiRepository) GetNotifications() (*dto.PbSubscriptionResponse, api_error.ApiErr) {
	var pbResp dto.PbSubscriptionResponse

	reqUrl, _ := url.Parse(r.cfg.PbApi.BaseUrl)
	reqUrl.Path = "/webhook"
	req, err := r.PrepareHttpRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	body, err := r.ExecHttpRequest(req)
	if err != nil {
		return nil, err
	}
	jsonErr := json.Unmarshal(*body, &pbResp)
	if jsonErr != nil {
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

func (r PbApiRepository) UnregisterForNotifications(notifs dto.PbSubscriptionResponse) api_error.ApiErr {
	for _, val := range notifs.Data {
		urlExt := fmt.Sprintf("/%v", val.ID)
		reqUrl, _ := url.Parse(r.cfg.PbApi.BaseUrl + "webhooks" + urlExt)
		req, err := r.PrepareHttpRequest("DELETE", reqUrl.String(), nil)
		if err != nil {
			return err
		}
		_, err = r.ExecHttpRequest(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r PbApiRepository) PrepareHttpRequest(reqType string, url string, body io.Reader) (*http.Request, api_error.ApiErr) {
	if reqType == "" {
		msg := "Request type cannot be empty"
		logger.Error(msg, nil)
		return nil, api_error.NewInternalServerError(msg, nil)
	}
	req, reqErr := http.NewRequest(reqType, url, body)
	if reqErr != nil {
		msg := "Could not create http request"
		logger.Error(msg, reqErr)
		return nil, api_error.NewInternalServerError(msg, reqErr)
	}
	authStr := fmt.Sprintf("Bearer %v", r.cfg.PbApi.ApiToken)
	req.Header = http.Header{
		"X-Version":     []string{"1"},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{authStr},
	}
	return req, nil
}

func (r PbApiRepository) ExecHttpRequest(req *http.Request) (*[]byte, api_error.ApiErr) {
	client := http.Client{}
	resp, resErr := client.Do(req)
	if resErr != nil {
		msg := "Error when executing http request"
		logger.Error(msg, resErr)
		return nil, api_error.NewInternalServerError(msg, resErr)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		msg := fmt.Sprintf("Error when sending request to Productboard API. Status code: %v. Message: %v", resp.StatusCode, string(body))
		logger.Error(msg, nil)
		return nil, api_error.NewInternalServerError(msg, nil)
	} else {
		logger.Info(fmt.Sprintf("Successfully sent request to Productboard API. Status code: %v", resp.StatusCode))
		return &body, nil
	}
}

func (r PbApiRepository) CreateSubscriptionRequest() (*[]byte, api_error.ApiErr) {
	subReq := dto.PbSubscriptionRequest{
		Data: dto.SubReqData{
			Name: "Feature Webhooks",
			Events: []dto.Events{
				{EventType: dto.PbEventTypes["featureCreate"]},
				{EventType: dto.PbEventTypes["featureUpdate"]},
				{EventType: dto.PbEventTypes["featureDelete"]},
			},
			Notification: dto.Notification{
				URL:     r.cfg.PbApi.WebHookUrl,
				Version: 1,
				Headers: dto.Headers{
					Authorization: r.cfg.RunTime.CallbackAuthToken,
				},
			},
		},
	}
	subReqJson, reqErr := json.Marshal(subReq)
	if reqErr != nil {
		msg := "Could not generate subscription request"
		logger.Error(msg, reqErr)
		return nil, api_error.NewInternalServerError(msg, reqErr)
	}
	return &subReqJson, nil
}
