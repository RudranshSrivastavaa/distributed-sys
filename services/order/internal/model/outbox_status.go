package model

type OutboxStatus string

const (
	OutboxPending OutboxStatus = "PENDING"

	OutboxPublished OutboxStatus = "PUBLISHED"

	OutboxFailed OutboxStatus = "FAILED"
)