package service

import (
	"context"
	stdErrors "errors"
	"log"

	//"log"

	//"time"
	"github.com/google/uuid"
	//ordercommand "github.com/rudransh/distributed-commerce/order/command"
	"github.com/rudransh/distributed-commerce/order/event"
	"github.com/rudransh/distributed-commerce/order/internal/dto"
	"github.com/rudransh/distributed-commerce/order/internal/errors"
	"github.com/rudransh/distributed-commerce/order/internal/mapper"
	"github.com/rudransh/distributed-commerce/order/internal/model"
	"github.com/rudransh/distributed-commerce/order/internal/outbox"
	"github.com/rudransh/distributed-commerce/order/internal/repository"
	"github.com/rudransh/distributed-commerce/order/internal/state"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	sagaevent "github.com/rudransh/distributed-commerce/saga/event"
	"gorm.io/gorm"
)

type orderService struct {
	repository repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{
		repository: repo,
	}
}

func (s *orderService) Create(request dto.CreateOrderRequest) (dto.OrderResponse, bool, error) {

	// 1. Check if this request has already been processed
	record, err := s.repository.FindByIdempotencyKey(request.IdempotencyKey)

	if err != nil {
		return dto.OrderResponse{}, false, err
	}

	if record != nil {
		// Idempotency key already exists
		order, err := s.repository.FindByID(record.OrderID)
		if err != nil {
			return dto.OrderResponse{}, false, err
		}

		return mapper.ToOrderResponse(order), false, nil
	}

	// Ignore "record not found", but return any other DB error
	if err != nil && !stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return dto.OrderResponse{}, false, err
	}

	// 2. Convert DTO -> Domain
	order := mapper.ToOrder(request)

	// 3. Validate
	if err := ValidateOrder(order); err != nil {
		return dto.OrderResponse{}, false, err
	}

	// 4. Apply business rules
	order.Status = model.StatusCreated
	CalculateTotal(order)

	// 5. Transaction
	err = s.repository.WithTransaction(func(tx *gorm.DB) error {

		txRepo := repository.NewOrderRepository(tx)

		outboxRepo := repository.NewOutboxRepository(tx)

		publisher := outbox.NewOutboxPublisher(outboxRepo)

		//----------------------------------------
		// Create Order
		//----------------------------------------

		if err := txRepo.Create(order); err != nil {
			return err
		}

		//----------------------------------------
		// Build Event
		//----------------------------------------

		evt, err := event.BuildOrderCreatedEvent(order)
		if err != nil {
			return err
		}

		//----------------------------------------
		// Save Outbox Event
		//----------------------------------------

		if err := publisher.Publish(
			context.Background(),
			kafkaa.OrderEvents,
			order.ID.String(),
			evt,
		); err != nil {
			return err
		}

		//----------------------------------------
		// Save Idempotency Key
		//----------------------------------------

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
		return dto.OrderResponse{}, false, err
	}

	// 6. Return response
	return mapper.ToOrderResponse(order), true, nil
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

func (s *orderService) GetByID(id uuid.UUID) (dto.OrderResponse, error) {

	order, err := s.repository.FindByID(id)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return mapper.ToOrderResponse(order), nil

}

func (s *orderService) Update(id uuid.UUID,request dto.UpdateOrderRequest) (dto.OrderResponse, error) {

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
func (s *orderService) Delete(id uuid.UUID) error {

	order, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if order.Status == model.StatusDelivered {
		return errors.ErrDeliveredOrderCannotBeDeleted
	}

	return s.repository.Delete(id)
}

func (s *orderService) UpdateStatus(id uuid.UUID,request dto.UpdateOrderStatusRequest) (dto.OrderResponse, error) {

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


func (s *orderService) HandleCompleteOrder(
	ctx context.Context,
	request sagaevent.CompleteOrderPayload,
) error {

	log.Println("Completing order:", request.OrderID)

	orderID, err := uuid.Parse(request.OrderID)
	if err != nil {
		return err
	}

	order, err := s.repository.FindByID(orderID)
	if err != nil {
		return err
	}

	order.Status = model.StatusPaid

	if err := s.repository.Update(order); err != nil {
		return err
	}

	log.Println("Order completed")

	return nil
}

func (s *orderService) HandleCancelOrderCommand(

	ctx context.Context,

	request sagaevent.CancelOrderPayload,

) error {

	log.Println("Cancelling order...")

	orderID, err := uuid.Parse(request.OrderID)
	if err != nil {
		return err
	}

	order, err := s.repository.FindByID(orderID)
	if err != nil {
		return err
	}

	order.Status = model.StatusCancelled

	if err := s.repository.Update(order); err != nil {
		return err
	}

	log.Println("Order cancelled")

	return nil
}