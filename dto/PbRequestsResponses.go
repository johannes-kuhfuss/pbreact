package dto

import "time"

var (
	PbEventTypes = map[string]string{
		"featureCreate": "feature.created",
		"featureUpdate": "feature.updated",
		"featureDelete": "feature.deleted",
	}
)

type PbSubscriptionRequest struct {
	Data SubReqData `json:"data"`
}

type PbSubscriptionResponse struct {
	Data  []SubRespData `json:"data"`
	Links Links         `json:"links"`
}

type EventNotification struct {
	Data EventData `json:"data"`
}

type Events struct {
	EventType string `json:"eventType"`
}

type Headers struct {
	Authorization string `json:"authorization"`
}
type Notification struct {
	URL     string  `json:"url"`
	Version int     `json:"version"`
	Headers Headers `json:"headers"`
}
type SubReqData struct {
	Name         string       `json:"name"`
	Events       []Events     `json:"events"`
	Notification Notification `json:"notification"`
}

type SubRespData struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Events    []Events  `json:"events"`
}
type Links struct {
	Next interface{} `json:"next"`
}

type EventData struct {
	ID        string     `json:"id"`
	EventType string     `json:"eventType"`
	Links     EventLinks `json:"links"`
}

type EventLinks struct {
	Target string `json:"target"`
}
