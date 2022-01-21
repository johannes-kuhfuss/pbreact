package dto

import "time"

type PbSubscriptionResponse struct {
	Data  []Data `json:"data"`
	Links Links  `json:"links"`
}
type Events struct {
	EventType string `json:"eventType"`
}
type Data struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Events    []Events  `json:"events"`
}
type Links struct {
	Next interface{} `json:"next"`
}
