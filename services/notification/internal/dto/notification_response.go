package dto

import "time"

type NotificationResponse struct {

	ID string `json:"id"`

	EventID string `json:"eventId"`

	Recipient string `json:"recipient"`

	Subject string `json:"subject"`

	Body string `json:"body"`

	Channel string `json:"channel"`

	Status string `json:"status"`

	FailureReason string `json:"failureReason"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`

}