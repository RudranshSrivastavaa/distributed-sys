package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
	sagaevent "github.com/rudransh/distributed-commerce/saga/event"
)

type InventoryService interface {
	CreateProduct(request dto.CreateProductRequest) (dto.ProductResponse, error)

	AddStock(productID uuid.UUID, request dto.AddStockRequest) (dto.InventoryResponse, error)

	Reserve(ctx context.Context, request dto.CreateReservationRequest) (dto.ReservationResponse, error)

	Release(ctx context.Context, reservationID uuid.UUID) (dto.ReservationResponse, error)

	Confirm(reservationID uuid.UUID) (dto.ReservationResponse, error)

	ExpireReservation(reservationID uuid.UUID) error

	GetExpiredReservations() ([]model.Reservation, error)

	ReserveInventory(
		ctx context.Context,
		request sagaevent.ReserveInventoryPayload,
	) error

	ReleaseInventory(
		ctx context.Context,
		request sagaevent.ReleaseInventoryPayload,
	) error
}
