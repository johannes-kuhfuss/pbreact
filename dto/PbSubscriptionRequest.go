package dto

type PbSubscriptionRequest struct {
	Data PbData `json:"data"`
}
type PbEvent struct {
	EventType string `json:"eventType"`
}
type PbHeaders struct {
	Authorization string `json:"authorization"`
}
type PbNotification struct {
	URL     string    `json:"url"`
	Version int       `json:"version"`
	Headers PbHeaders `json:"headers"`
}
type PbData struct {
	Name         string         `json:"name"`
	Events       []PbEvent      `json:"events"`
	Notification PbNotification `json:"notification"`
}
