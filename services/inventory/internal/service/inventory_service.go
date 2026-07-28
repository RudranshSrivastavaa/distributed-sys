package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/inventory/event"
	"github.com/rudransh/distributed-commerce/inventory/internal/database"
	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	inventory_error "github.com/rudransh/distributed-commerce/inventory/internal/errors"
	"github.com/rudransh/distributed-commerce/inventory/internal/mapper"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
	"github.com/rudransh/distributed-commerce/inventory/internal/outbox"
	"github.com/rudransh/distributed-commerce/inventory/internal/repository"
	"github.com/rudransh/distributed-commerce/inventory/internal/state"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	sagaevent "github.com/rudransh/distributed-commerce/saga/event"

	"gorm.io/gorm"
)

type inventoryService struct {
	productRepository     repository.ProductRepository
	inventoryRepository   repository.InventoryRepository
	reservationRepository repository.ReservationRepository
	txManager             *repository.TransactionManager
}

func NewInventoryService(
	productRepo repository.ProductRepository,
	inventoryRepo repository.InventoryRepository,
	reservationRepo repository.ReservationRepository,
	txManager *repository.TransactionManager,
) InventoryService {

	return &inventoryService{
		productRepository:     productRepo,
		inventoryRepository:   inventoryRepo,
		reservationRepository: reservationRepo,
		txManager:             txManager,
	}
}

func (s *inventoryService) CreateProduct(
	request dto.CreateProductRequest,
) (dto.ProductResponse, error) {

	if err := validateCreateProduct(request); err != nil {
		return dto.ProductResponse{}, err
	}

	if _, err := s.productRepository.FindBySKU(request.SKU); err == nil {
		return dto.ProductResponse{}, errors.New("sku already exists")
	}

	product := mapper.ToProduct(request)

	err := s.txManager.Execute(func(tx *gorm.DB) error {

		productRepo := repository.NewProductRepository(tx)
		inventoryRepo := repository.NewInventoryRepository(tx)

		if err := productRepo.Create(product); err != nil {
			return err
		}

		inventory := &model.Inventory{
			ProductID:         product.ID,
			AvailableQuantity: 0,
			ReservedQuantity:  0,
			Version:           1,
		}

		if err := inventoryRepo.Create(inventory); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.ProductResponse{}, err
	}

	return mapper.ToProductResponse(product), nil
}

func (s *inventoryService) AddStock(productID uuid.UUID, request dto.AddStockRequest) (dto.InventoryResponse, error) {

	if err := validateStockQuantity(
		request.Quantity,
	); err != nil {

		return dto.InventoryResponse{}, err
	}

	inventory, err := s.inventoryRepository.
		FindByProductID(productID)

	if err != nil {
		return dto.InventoryResponse{}, err
	}

	inventory.AvailableQuantity += request.Quantity

	if err := s.inventoryRepository.
		Update(inventory); err != nil {

		return dto.InventoryResponse{}, err
	}

	return mapper.ToInventoryResponse(
		inventory,
	), nil

}

func (s *inventoryService) Reserve(ctx context.Context,request dto.CreateReservationRequest) (dto.ReservationResponse, error) {

	if request.Quantity <= 0 {
		return dto.ReservationResponse{},
			errors.New("quantity must be greater than zero")
	}

	var reservation *model.Reservation

	const maxRetries = 3

	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {

		err = s.txManager.Execute(func(tx *gorm.DB) error {

			inventoryRepo := repository.NewInventoryRepository(tx)
			reservationRepo := repository.NewReservationRepository(tx)

			outboxRepo := repository.NewOutboxRepository(tx)
			publisher := outbox.NewOutboxPublisher(outboxRepo)

			log.Println("1. Find inventory")

			inventory, err := inventoryRepo.FindByProductID(
				request.ProductID,
			)
			if err != nil {
				return err
			}
			log.Println("2. Reserve inventory")
			if err := inventory.Reserve(request.Quantity); err != nil {
				return err
			}
			log.Println("3. Update inventory")
			if err := inventoryRepo.UpdateWithVersion(inventory); err != nil {
				return err
			}

			reservation = &model.Reservation{
				OrderID:   request.OrderID,
				ProductID: request.ProductID,
				Quantity:  request.Quantity,
				Status:    model.StatusReserved,
				ExpiresAt: time.Now().Add(15 * time.Minute),
			}
			log.Println("4. Create reservation")
			if err := reservationRepo.Create(reservation); err != nil {
				return err
			}
			log.Println("5. Write outbox")
			evt, err := event.BuildInventoryReservedEvent(
				reservation,
			)

			if err != nil {
				return err
			}

			if err := publisher.Publish(
				ctx,
				kafkaa.InventoryEvents,
				reservation.ID.String(),
				evt,
			); err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			log.Println("6. Success")
			break
		}

		if !errors.Is(err, inventory_error.ErrOptimisticLockConflict) {
			break
		}
	}

	if err != nil {
		return dto.ReservationResponse{}, err
	}

	return mapper.ToReservationResponse(
		reservation,
	), nil

}

func (s *inventoryService) Release(ctx context.Context,reservationID uuid.UUID) (dto.ReservationResponse, error) {

	var reservation *model.Reservation

	err := s.txManager.Execute(func(tx *gorm.DB) error {

		reservationRepo := repository.NewReservationRepository(tx)

		inventoryRepo := repository.NewInventoryRepository(tx)

		outboxRepo := repository.NewOutboxRepository(tx)

		publisher := outbox.NewOutboxPublisher(outboxRepo)

		var err error

		reservation, err = reservationRepo.FindByID(reservationID)

		if err != nil {
			return err
		}

		inventory, err := inventoryRepo.FindByProductID(
			reservation.ProductID,
		)

		if err != nil {
			return err
		}

		if err := inventory.Release(reservation.Quantity); err != nil {
			return err
		}

		if err := state.Transition(

			reservation,

			model.StatusReleased,
		); err != nil {

			return err

		}

		if err := inventoryRepo.UpdateWithVersion(inventory); err != nil {
			return err
		}

		if err := reservationRepo.Update(reservation); err != nil {
			return err
		}

		evt, err := event.BuildInventoryReleasedEvent(
			reservation,
		)

		if err != nil {
			return err
		}

		if err := publisher.Publish(
			ctx,
			kafkaa.InventoryEvents,
			reservation.ID.String(),
			evt,
		); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return dto.ReservationResponse{}, err
	}

	return mapper.ToReservationResponse(reservation), nil
}

func (s *inventoryService) Confirm(reservationID uuid.UUID) (dto.ReservationResponse, error) {

	var reservation *model.Reservation

	err := s.txManager.Execute(func(tx *gorm.DB) error {

		reservationRepo := repository.NewReservationRepository(tx)

		inventoryRepo := repository.NewInventoryRepository(tx)

		var err error

		reservation, err = reservationRepo.FindByID(reservationID)

		if err != nil {
			return err
		}

		inventory, err := inventoryRepo.FindByProductID(
			reservation.ProductID,
		)

		if err != nil {
			return err
		}

		if err := inventory.Confirm(
			reservation.Quantity,
		); err != nil {
			return err
		}

		if err := state.Transition(

			reservation,

			model.StatusConfirmed,
		); err != nil {

			return err

		}

		if err := inventoryRepo.UpdateWithVersion(inventory); err != nil {
			return err
		}

		if err := reservationRepo.Update(reservation); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.ReservationResponse{}, err
	}

	return mapper.ToReservationResponse(reservation), nil
}

func (s *inventoryService) GetExpiredReservations() ([]model.Reservation,error) {

	return s.reservationRepository.
		FindExpiredReservations(
			time.Now(),
		)

}

func (s *inventoryService) ExpireReservation(
	reservationID uuid.UUID,
) error {

	return s.txManager.Execute(func(tx *gorm.DB) error {

		reservationRepo := repository.NewReservationRepository(tx)

		inventoryRepo := repository.NewInventoryRepository(tx)

		reservation, err := reservationRepo.FindByID(reservationID)
		if err != nil {
			return err
		}

		if reservation.Status != model.StatusReserved {
			return nil
		}

		inventory, err := inventoryRepo.FindByProductID(
			reservation.ProductID,
		)
		if err != nil {
			return err
		}

		if err := inventory.Release(
			reservation.Quantity,
		); err != nil {
			return err
		}

		if err := state.Transition(
			reservation,
			model.StatusExpired,
		); err != nil {
			return err
		}

		if err := inventoryRepo.UpdateWithVersion(
			inventory,
		); err != nil {
			return err
		}

		if err := reservationRepo.Update(
			reservation,
		); err != nil {
			return err
		}

		return nil
	})
}

func (s *inventoryService) ReserveInventory(
	ctx context.Context,
	request sagaevent.ReserveInventoryPayload,
) error {

	for _, item := range request.Items {

		_, err := s.Reserve(ctx,

			dto.CreateReservationRequest{

				OrderID: request.OrderID,

				ProductID: item.ProductID,

				Quantity: item.Quantity,
			},
		)

		if err != nil {

			log.Println("Inventory reservation failed")

			evt, buildErr := event.BuildInventoryReservationFailedEvent(

				request.OrderID.String(),

				item.ProductID.String(),

				item.Quantity,

				err.Error(),
			)

			if buildErr != nil {
				return buildErr
			}
			outboxRepo := repository.NewOutboxRepository(database.DB)

			publisher := outbox.NewOutboxPublisher(outboxRepo)

			if err := publisher.Publish(

				ctx,

				kafkaa.InventoryEvents,

				request.OrderID.String(),

				evt,
			); err != nil {

				return err
			}

			log.Println("Published INVENTORY_RESERVATION_FAILED")

			return nil
		}
	}
	return nil
}


func (s *inventoryService) ReleaseInventory(

	ctx context.Context,

	request sagaevent.ReleaseInventoryPayload,

) error {

	log.Println("Releasing inventory...")

	reservations, err := s.reservationRepository.FindByOrderIDD(
		request.OrderID,
	)

	if err != nil {
		return err
	}

	for _, reservation := range reservations {

		_, err := s.Release(
			ctx,
			reservation.ID,
		)

		if err != nil {
			return err
		}
	}

	log.Println("Inventory released")

	return nil
}