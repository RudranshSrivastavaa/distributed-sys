package model

type SagaStatus string

const (

	SagaRunning SagaStatus = "RUNNING"

	SagaCompleted SagaStatus = "COMPLETED"

	SagaCompensating SagaStatus = "COMPENSATING"

	SagaFailed SagaStatus = "FAILED"
)