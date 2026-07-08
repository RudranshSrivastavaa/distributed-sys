package model

type PaymentStatus string

const (
	StatusCreated PaymentStatus = "CREATED"

	StatusPending PaymentStatus = "PENDING"

	StatusSuccess PaymentStatus = "SUCCESS"

	StatusFailed PaymentStatus = "FAILED"

	StatusRefunded PaymentStatus = "REFUNDED"
)

func (s PaymentStatus) IsFinal() bool {

	return s == StatusSuccess ||
		s == StatusFailed ||
		s == StatusRefunded

}

func (s PaymentStatus) CanProcess() bool {

	return s == StatusCreated ||
		s == StatusPending

}