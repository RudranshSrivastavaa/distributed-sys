package service

import (
	"github.com/rudransh/distributed-commerce/payment/internal/model"
)



func (s *paymentService) UpdateStatus(
	payment *model.Payment,
	status model.PaymentStatus,
) {

	payment.Status = status

}
