package mapper

import (
	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

func ToReservationResponse(
	r *model.Reservation,
) dto.ReservationResponse {

	return dto.ReservationResponse{
		ID: r.ID,

		OrderID: r.OrderID,

		ProductID: r.ProductID,

		Quantity: r.Quantity,

		Status: r.Status,

		ExpiresAt: r.ExpiresAt,
	}

}