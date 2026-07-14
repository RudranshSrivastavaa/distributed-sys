package model

type StepStatus string

const (

	StepPending StepStatus = "PENDING"

	StepCompleted StepStatus = "COMPLETED"

	StepFailed StepStatus = "FAILED"
)