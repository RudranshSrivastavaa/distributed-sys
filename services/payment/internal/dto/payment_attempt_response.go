package dto

import (
	"time"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type PaymentAttemptResponse struct {

	AttemptNumber int `json:"attemptNumber"`

	Status model.PaymentStatus `json:"status"`

	FailureReason string `json:"failureReason,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
}