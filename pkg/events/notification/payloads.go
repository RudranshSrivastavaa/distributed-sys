package notification

type NotificationSentPayload struct {

	NotificationID string `json:"notificationId"`

	EventID string `json:"eventId"`

	Recipient string `json:"recipient"`
}

type NotificationFailedPayload struct {

	NotificationID string `json:"notificationId"`

	EventID string `json:"eventId"`

	Reason string `json:"reason"`
}