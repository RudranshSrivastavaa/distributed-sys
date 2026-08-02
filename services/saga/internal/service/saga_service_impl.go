package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	sagaevent "github.com/rudransh/distributed-commerce/saga/event"
	"github.com/rudransh/distributed-commerce/saga/internal/model"
	"github.com/rudransh/distributed-commerce/saga/internal/repository"

	inventoryevents "github.com/rudransh/distributed-commerce/pkg/events/inventory"
	orderevents "github.com/rudransh/distributed-commerce/pkg/events/order"
	paymentevents "github.com/rudransh/distributed-commerce/pkg/events/payment"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type sagaService struct {
	repository repository.SagaRepository
	producer kafkaa.Producer
}

func NewSagaService(repository repository.SagaRepository,producer kafkaa.Producer) SagaService {

	return &sagaService{

		repository: repository,
		producer: producer,
	}
}

func (s *sagaService) StartSaga(ctx context.Context,payload orderevents.OrderCreatedPayload,
) error {

	log.Println("Starting saga for order:", payload.OrderID)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return err
	}

	//------------------------------------------
	// Create Saga
	//------------------------------------------

	saga := &model.Saga{

		OrderID: payload.OrderID,

		Amount: payload.TotalPrice,

		Currency: payload.Currency,

		Status: model.SagaRunning,

		InventoryStatus: model.StepPending,

		PaymentStatus: model.StepPending,
	}

	if err := s.repository.Create(saga); err != nil {
		return err
	}
	log.Println("Saga created")

	//------------------------------------------
	// Build Command
	//------------------------------------------

	items := make([]sagaevent.ReserveInventoryItem, 0, len(payload.Items))

	for _, item := range payload.Items {

		items = append(items, sagaevent.ReserveInventoryItem{

			ProductID: item.ProductID,

			Quantity: int64(item.Quantity),
		})
	}

	command, err := sagaevent.BuildReserveInventoryCommand(orderID,payload.CustomerID,items)

	if err != nil {
		return err
	}

	//------------------------------------------
	// Publish Command
	//------------------------------------------
	log.Println("Publishing RESERVE_INVENTORY")
	return s.producer.Publish(

		ctx,

		kafkaa.SagaCommands.Name,

		payload.OrderID,

		command,
	)
}

func (s *sagaService) HandleInventoryReserved(ctx context.Context,payload inventoryevents.InventoryReservedPayload) error {

	log.Println("Updating saga...")

	saga, err := s.repository.FindByOrderID(
		payload.OrderID,
	)

	if err != nil {
		return err
	}

	saga.InventoryStatus = model.StepCompleted

	if err := s.repository.Update(saga); err != nil {
		return err
	}

	log.Println("Inventory step completed")

	//------------------------------------------
	// Build PROCESS_PAYMENT command
	//------------------------------------------

	sagaOrderID, err := uuid.Parse(saga.OrderID)
	if err != nil {
		return err
	}

	command, err := sagaevent.BuildProcessPaymentCommand(
		sagaOrderID,
		saga.Amount,
		saga.Currency,
	)

	if err != nil {
		return err
	}

	log.Println("Publishing PROCESS_PAYMENT")

	err = s.producer.Publish(
		ctx,
		kafkaa.SagaCommands.Name,
		payload.OrderID,
		command,
	)

	if err != nil {
		log.Println("Publish failed:", err)
		return err
	}

	log.Println("PROCESS_PAYMENT published successfully")

	return nil
}

func (s *sagaService) HandleInventoryReservationFailed(

	ctx context.Context,

	payload inventoryevents.InventoryReservationFailedPayload,

) error {

	log.Println("Inventory reservation failed")

	saga, err := s.repository.FindByOrderID(
		payload.OrderID,
	)

	if err != nil {
		return err
	}

	saga.InventoryStatus = model.StepFailed

	saga.Status = model.SagaFailed

	if err := s.repository.Update(
		saga,
	); err != nil {
		return err
	}

	log.Println("Saga updated")

	command, err := sagaevent.BuildCancelOrderCommand(
		payload.OrderID,
	)

	if err != nil {
		return err
	}

	log.Println("Publishing CANCEL_ORDER")

	return s.producer.Publish(

		ctx,

		kafkaa.SagaCommands.Name,

		payload.OrderID,

		command,
	)
}

func (s *sagaService) HandlePaymentCompleted(ctx context.Context,payload paymentevents.PaymentSucceededPayload) error {

	log.Println("1. Finding saga")

	saga, err := s.repository.FindByOrderID(payload.OrderID)
	if err != nil {
		log.Println("FindByOrderID failed:", err)
		return err
	}

	log.Println("2. Saga found")

	saga.PaymentStatus = model.StepCompleted
	saga.Status = model.SagaCompleted

	log.Println("3. Updating saga")

	if err := s.repository.Update(saga); err != nil {
		log.Println("Update failed:", err)
		return err
	}

	log.Println("4. Saga updated")

	command, err := sagaevent.BuildCompleteOrderCommand(
		payload.OrderID,
	)
	if err != nil {
		log.Println("BuildCompleteOrderCommand failed:", err)
		return err
	}

	log.Println("5. Command built")

	err = s.producer.Publish(
		ctx,
		kafkaa.SagaCommands.Name,
		payload.OrderID,
		command,
	)
	if err != nil {
		log.Println("Publish failed:", err)
		return err
	}

	log.Println("6. COMPLETE_ORDER published")

	return nil
}


func (s *sagaService) HandlePaymentFailed(ctx context.Context,payload paymentevents.PaymentFailedPayload) error {

	log.Println("Updating saga...")

	saga, err := s.repository.FindByOrderID(
		payload.OrderID,
	)
	if err != nil {
		return err
	}

	saga.PaymentStatus = model.StepFailed
	saga.Status = model.SagaFailed

	if err := s.repository.Update(saga); err != nil {
		return err
	}

	log.Println("Payment step failed")

	command, err := sagaevent.BuildReleaseInventoryCommand(
		payload.OrderID,
	)

	if err != nil {
		return err
	}

	log.Println("Publishing RELEASE_INVENTORY")

	return s.producer.Publish(

		ctx,

		kafkaa.SagaCommands.Name,

		payload.OrderID,

		command,
	)
}

func (s *sagaService) HandleInventoryReleased(
	ctx context.Context,
	payload inventoryevents.InventoryReleasedPayload,
) error {

	log.Println("Updating saga...")

	saga, err := s.repository.FindByOrderID(payload.OrderID)
	if err != nil {
		return err
	}

	saga.InventoryStatus = model.StepFailed
	saga.Status = model.SagaCompensating


	if err := s.repository.Update(saga); err != nil {
		return err
	}

	log.Println("Inventory compensation completed")

	command, err := sagaevent.BuildCancelOrderCommand(
		payload.OrderID,
	)
	if err != nil {
		return err
	}

	log.Println("Publishing CANCEL_ORDER")

	return s.producer.Publish(
		ctx,
		kafkaa.SagaCommands.Name,
		payload.OrderID,
		command,
	)
}