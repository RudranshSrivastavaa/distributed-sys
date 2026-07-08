package provider

import "github.com/rudransh/distributed-commerce/payment/internal/model"

type PaymentResult struct {

	Success bool
	
	Status model.PaymentStatus

	Reference string

	Message string

}