package model

type RefundStatus string

const (

	RefundStatusPending RefundStatus = "PENDING"

	RefundStatusSucceeded RefundStatus = "SUCCEEDED"

	RefundStatusFailed RefundStatus = "FAILED"
)

func (s RefundStatus) IsFinal() bool {

	return s == RefundStatusSucceeded ||
		s == RefundStatusFailed
}

func (s RefundStatus) CanProcess() bool {

	return s == RefundStatusPending
}