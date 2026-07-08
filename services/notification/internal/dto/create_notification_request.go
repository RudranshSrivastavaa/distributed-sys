package dto

type CreateNotificationRequest struct {

	EventID string `json:"eventId"`

	Recipient string `json:"recipient"`

	Subject string `json:"subject"`

	Body string `json:"body"`

	Channel string `json:"channel"`

}