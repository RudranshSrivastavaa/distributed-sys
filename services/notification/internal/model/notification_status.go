package model

type NotificationStatus string

const (

	StatusPending NotificationStatus = "PENDING"

	StatusProcessing NotificationStatus = "PROCESSING"

	StatusSent NotificationStatus = "SENT"

	StatusFailed NotificationStatus = "FAILED"

)

func (s NotificationStatus) IsFinal() bool {

	return s == StatusSent || s == StatusFailed

}

func (s NotificationStatus) CanProcess() bool {

	return s == StatusPending || s == StatusFailed

}