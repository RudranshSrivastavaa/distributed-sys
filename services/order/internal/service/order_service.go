package service

import (
	stdErrors "errors"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/order/internal/dto"
	"github.com/rudransh/distributed-commerce/order/internal/errors"
	"github.com/rudransh/distributed-commerce/order/internal/mapper"
	"github.com/rudransh/distributed-commerce/order/internal/model"
	"github.com/rudransh/distributed-commerce/order/internal/repository"
	"github.com/rudransh/distributed-commerce/order/internal/state"
	"github.com/rudransh/distributed-commerce/pkg/kafka"
	"gorm.io/gorm"
)

type orderService struct {
	repository repository.OrderRepository
	producer kafka.Producer
}

func NewOrderService(
	repo repository.OrderRepository,
	producer kafka.Producer,
) OrderService {

	return &orderService{
		repository: repo,
		producer: producer,
	}
}

func (s *orderService) Create(request dto.CreateOrderRequest) (dto.OrderResponse ,bool , error) {

	// 1. Check if this request has already been processed
	record, err := s.repository.FindByIdempotencyKey(request.IdempotencyKey)

	if err != nil {
	return dto.OrderResponse{}, false ,err
    } 

	if record!= nil {
		// Idempotency key already exists
		order, err := s.repository.FindByID(record.OrderID)
		if err != nil {
			return dto.OrderResponse{},false , err
		}

		return mapper.ToOrderResponse(order),false , nil
	}

	// Ignore "record not found", but return any other DB error
	if err != nil && !stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return dto.OrderResponse{},false , err
	}

	// 2. Convert DTO -> Domain
	order := mapper.ToOrder(request)

	// 3. Validate
	if err := ValidateOrder(order); err != nil {
		return dto.OrderResponse{},false , err
	}

	// 4. Apply business rules
	order.Status = model.StatusCreated
	CalculateTotal(order)

	// 5. Transaction
	err = s.repository.WithTransaction(func(tx *gorm.DB) error {

		txRepo := repository.NewOrderRepository(tx)

		// Create Order
		if err := txRepo.Create(order); err != nil {
			return err
		}

		// Save Idempotency Key
		record := &model.IdempotencyKey{
			ID:      uuid.New(),
			Key:     request.IdempotencyKey,
			OrderID: order.ID,
		}

		if err := txRepo.SaveIdempotencyKey(record); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.OrderResponse{},false , err
	}

	// 6. Return response
	return mapper.ToOrderResponse(order), true ,nil
}

func (s *orderService) GetAll() ([]dto.OrderResponse, error) {

	orders, err := s.repository.FindAll()

	if err != nil {
		return nil, err
	}

	return mapper.ToOrderResponses(
		orders,
	), nil

}

func (s *orderService) GetByID(
	id uuid.UUID,
) (dto.OrderResponse, error) {

	order, err := s.repository.FindByID(id)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return mapper.ToOrderResponse(order), nil

}

func (s *orderService) Update(
	id uuid.UUID,
	request dto.UpdateOrderRequest,
) (dto.OrderResponse, error) {

	order, err := s.repository.FindByID(id)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	mapper.UpdateOrder(order, request)

	if err := ValidateOrder(order); err != nil {
		return dto.OrderResponse{}, err
	}

	CalculateTotal(order)

	if err := s.repository.Update(order); err != nil {
		return dto.OrderResponse{}, err
	}

	return mapper.ToOrderResponse(order), nil
}
func (s *orderService) Delete(
	id uuid.UUID,
) error {

	order, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if order.Status == model.StatusDelivered {
		return errors.ErrDeliveredOrderCannotBeDeleted
	}

	return s.repository.Delete(id)
}

func (s *orderService) UpdateStatus(
	id uuid.UUID,
	request dto.UpdateOrderStatusRequest,
) (dto.OrderResponse, error) {

	order, err := s.repository.FindByID(id)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	if err := state.Transition(
		order,
		request.Status,
	); err != nil {
		return dto.OrderResponse{}, err
	}

	if err := s.repository.Update(order); err != nil {
		return dto.OrderResponse{}, err
	}

	return mapper.ToOrderResponse(order), nil
}
