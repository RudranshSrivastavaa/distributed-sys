package service

import (
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

func (s *paymentService) toPaymentResponse(
	payment *model.Payment,
	attempts []model.PaymentAttempt,
	paymentURL string,
) dto.PaymentResponse {

	response := dto.PaymentResponse{
		ID:                payment.ID,
		OrderID:           payment.OrderID,
		Amount:            payment.Money.Amount,
		Currency:          payment.Money.Currency,
		Status:            payment.Status,
		Provider:          payment.Provider,
		ProviderReference: payment.ProviderReference,
		PaymentURL:        paymentURL,
		CreatedAt:         payment.CreatedAt,
		UpdatedAt:         payment.UpdatedAt,
	}

	for _, attempt := range attempts {

		response.Attempts = append(
			response.Attempts,
			dto.PaymentAttemptResponse{
				AttemptNumber: attempt.AttemptNumber,
				Status:        attempt.Status,
				FailureReason: attempt.FailureReason,
				CreatedAt:     attempt.CreatedAt,
			},
		)

	}

	return response

}