package events

type NotificationSentEvent struct {

	Metadata

	NotificationID string `json:"notificationId"`

	EventID string `json:"eventId"`

	Recipient string `json:"recipient"`
}


type NotificationFailedEvent struct {

	Metadata

	NotificationID string `json:"notificationId"`

	Reason string `json:"reason"`
}