package service

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type InventoryService interface {

	CreateProduct(request dto.CreateProductRequest) (dto.ProductResponse, error)

	AddStock(productID uuid.UUID,request dto.AddStockRequest) (dto.InventoryResponse, error)

	Reserve(request dto.CreateReservationRequest)(dto.ReservationResponse, error)

	Release(reservationID uuid.UUID) (dto.ReservationResponse, error)

	Confirm(reservationID uuid.UUID) (dto.ReservationResponse, error)

	ExpireReservation(reservationID uuid.UUID) error

	GetExpiredReservations() ([]model.Reservation, error)


}